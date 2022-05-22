package provider

import (
	"context"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_project_mirror", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`" + `gitlab_project_mirror` + "`" + ` resource allows to manage the lifecycle of a project mirror.

This is for *pushing* changes to a remote repository. *Pull Mirroring* can be configured using a combination of the
import_url, mirror, and mirror_trigger_builds properties on the gitlab_project resource.

-> **Destroy Behavior** GitLab 14.10 introduced an API endpoint to delete a project mirror.
   Therefore, for GitLab 14.10 and newer the project mirror will be destroyed when the resource is destroyed.
   For older versions, the mirror will be disabled and the resource will be destroyed.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/remote_mirrors.html)`,

		CreateContext: resourceGitlabProjectMirrorCreate,
		ReadContext:   resourceGitlabProjectMirrorRead,
		UpdateContext: resourceGitlabProjectMirrorUpdate,
		DeleteContext: resourceGitlabProjectMirrorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Description: "The id of the project.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
			},
			"mirror_id": {
				Description: "Mirror ID.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"url": {
				Description: "The URL of the remote repository to be mirrored.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Sensitive:   true, // Username and password must be provided in the URL for https.
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					oldURL, err := url.Parse(old)
					if err != nil {
						return old == new
					}
					newURL, err := url.Parse(new)
					if err != nil {
						return old == new
					}
					if oldURL.User != nil {
						oldURL.User = url.UserPassword("redacted", "redacted")
					}
					if newURL.User != nil {
						newURL.User = url.UserPassword("redacted", "redacted")
					}
					return oldURL.String() == newURL.String()
				},
			},
			"enabled": {
				Description: "Determines if the mirror is enabled.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"only_protected_branches": {
				Description: "Determines if only protected branches are mirrored.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"keep_divergent_refs": {
				Description: "Determines if divergent refs are skipped.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
		},
	}
})

func resourceGitlabProjectMirrorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	projectID := d.Get("project").(string)
	URL := d.Get("url").(string)
	enabled := d.Get("enabled").(bool)
	onlyProtectedBranches := d.Get("only_protected_branches").(bool)
	keepDivergentRefs := d.Get("keep_divergent_refs").(bool)

	options := &gitlab.AddProjectMirrorOptions{
		URL:                   &URL,
		Enabled:               &enabled,
		OnlyProtectedBranches: &onlyProtectedBranches,
		KeepDivergentRefs:     &keepDivergentRefs,
	}

	log.Printf("[DEBUG] create gitlab project mirror for project %v", projectID)

	mirror, _, err := client.ProjectMirrors.AddProjectMirror(projectID, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("mirror_id", mirror.ID)

	mirrorID := strconv.Itoa(mirror.ID)
	d.SetId(buildTwoPartID(&projectID, &mirrorID))
	return resourceGitlabProjectMirrorRead(ctx, d, meta)
}

func resourceGitlabProjectMirrorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	mirrorID := d.Get("mirror_id").(int)
	projectID := d.Get("project").(string)
	enabled := d.Get("enabled").(bool)
	onlyProtectedBranches := d.Get("only_protected_branches").(bool)
	keepDivergentRefs := d.Get("keep_divergent_refs").(bool)

	options := gitlab.EditProjectMirrorOptions{
		Enabled:               &enabled,
		OnlyProtectedBranches: &onlyProtectedBranches,
		KeepDivergentRefs:     &keepDivergentRefs,
	}
	log.Printf("[DEBUG] update gitlab project mirror %v for %s", mirrorID, projectID)

	_, _, err := client.ProjectMirrors.EditProjectMirror(projectID, mirrorID, &options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceGitlabProjectMirrorRead(ctx, d, meta)
}

func resourceGitlabProjectMirrorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	mirrorID := d.Get("mirror_id").(int)
	projectID := d.Get("project").(string)

	isDeleteSupported, err := isGitLabVersionAtLeast(ctx, client, "14.10")()
	if err != nil {
		return diag.FromErr(err)
	}

	if isDeleteSupported {
		log.Printf("[DEBUG] delete gitlab project mirror %v for %s", mirrorID, projectID)

		_, err := client.ProjectMirrors.DeleteProjectMirror(projectID, mirrorID, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		// NOTE: this code only exists to support GitLab < 14.10.
		//       It can be removed once ~ GitLab 15.2 is out and supported.
		options := gitlab.EditProjectMirrorOptions{Enabled: gitlab.Bool(false)}
		log.Printf("[DEBUG] Disable gitlab project mirror %v for %s", mirrorID, projectID)
		_, _, err := client.ProjectMirrors.EditProjectMirror(projectID, mirrorID, &options, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceGitlabProjectMirrorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	ids := strings.Split(d.Id(), ":")
	projectID := ids[0]
	rawMirrorID := ids[1]
	mirrorID, err := strconv.Atoi(rawMirrorID)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] read gitlab project mirror %s id %v", projectID, mirrorID)
	mirror, err := resourceGitLabProjectMirrorGetMirror(ctx, client, projectID, mirrorID)
	if err != nil {
		return diag.FromErr(err)
	}

	if mirror == nil {
		log.Printf("[DEBUG] mirror %d in project %s not found, removing from state", mirrorID, projectID)
		d.SetId("")
		return nil
	}

	resourceGitlabProjectMirrorSetToState(d, mirror, &projectID)
	return nil
}

func resourceGitlabProjectMirrorSetToState(d *schema.ResourceData, projectMirror *gitlab.ProjectMirror, projectID *string) {
	d.Set("enabled", projectMirror.Enabled)
	d.Set("mirror_id", projectMirror.ID)
	d.Set("keep_divergent_refs", projectMirror.KeepDivergentRefs)
	d.Set("only_protected_branches", projectMirror.OnlyProtectedBranches)
	d.Set("project", projectID)
	d.Set("url", projectMirror.URL)
}

func resourceGitLabProjectMirrorGetMirror(ctx context.Context, client *gitlab.Client, projectID string, mirrorID int) (*gitlab.ProjectMirror, error) {
	isGetProjectMirrorSupported, err := isGitLabVersionAtLeast(ctx, client, "14.10")()
	if err != nil {
		return nil, err
	}

	var mirror *gitlab.ProjectMirror

	if isGetProjectMirrorSupported {
		mirror, _, err = client.ProjectMirrors.GetProjectMirror(projectID, mirrorID, gitlab.WithContext(ctx))
		if err != nil {
			if is404(err) {
				return nil, nil
			}
			return nil, err
		}
	} else {
		// NOTE: remove this branch and move logic back to Read() function when GitLab older than 14.10 are not longer supported by this provider
		found := false
		options := &gitlab.ListProjectMirrorOptions{
			Page:    1,
			PerPage: 20,
		}

		for options.Page != 0 && !found {
			mirrors, resp, err := client.ProjectMirrors.ListProjectMirror(projectID, options, gitlab.WithContext(ctx))
			if err != nil {
				return nil, err
			}

			for _, m := range mirrors {
				if m.ID == mirrorID {
					mirror = m
					found = true
					break
				}
			}
			options.Page = resp.NextPage
		}
	}

	return mirror, nil
}
