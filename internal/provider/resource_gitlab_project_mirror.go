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

func resourceGitlabProjectMirror() *schema.Resource {
	return &schema.Resource{
		Description: "This resource allows you to add a mirror target for the repository, all changes will be synced to the remote target.\n\n" +
			"-> This is for *pushing* changes to a remote repository. *Pull Mirroring* can be configured using a combination of the\n" +
			"`import_url`, `mirror`, and `mirror_trigger_builds` properties on the `gitlab_project` resource.\n\n" +
			"For further information on mirroring, consult the\n" +
			"[gitlab documentation](https://docs.gitlab.com/ee/user/project/repository/repository_mirroring.html#repository-mirroring).",

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
}

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

// Documented remote mirrors API does not support a delete method, instead mirror is disabled.
func resourceGitlabProjectMirrorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	enabled := false

	mirrorID := d.Get("mirror_id").(int)
	projectID := d.Get("project").(string)
	onlyProtectedBranches := d.Get("only_protected_branches").(bool)
	keepDivergentRefs := d.Get("keep_divergent_refs").(bool)

	options := gitlab.EditProjectMirrorOptions{
		Enabled:               &enabled,
		OnlyProtectedBranches: &onlyProtectedBranches,
		KeepDivergentRefs:     &keepDivergentRefs,
	}
	log.Printf("[DEBUG] Disable gitlab project mirror %v for %s", mirrorID, projectID)

	_, _, err := client.ProjectMirrors.EditProjectMirror(projectID, mirrorID, &options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGitlabProjectMirrorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	ids := strings.Split(d.Id(), ":")
	projectID := ids[0]
	mirrorID := ids[1]
	integerMirrorID, err := strconv.Atoi(mirrorID)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] read gitlab project mirror %s id %v", projectID, mirrorID)

	var mirror *gitlab.ProjectMirror
	found := false

	opts := &gitlab.ListProjectMirrorOptions{
		Page:    1,
		PerPage: 20,
	}

	for {
		mirrors, response, err := client.ProjectMirrors.ListProjectMirror(projectID, opts, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}

		for _, m := range mirrors {
			log.Printf("[DEBUG] project mirror found %v", m.ID)
			if m.ID == integerMirrorID {
				mirror = m
				found = true
				break
			}
		}
		if response.CurrentPage >= response.TotalPages {
			break
		}
		opts.Page++
	}

	if !found {
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
