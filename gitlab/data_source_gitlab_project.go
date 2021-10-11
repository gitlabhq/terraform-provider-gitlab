package gitlab

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func dataSourceGitlabProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGitlabProjectRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path_with_namespace": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_branch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_access_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"issues_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"merge_requests_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"pipelines_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"wiki_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"snippets_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"lfs_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"visibility_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ssh_url_to_repo": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"http_url_to_repo": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"web_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"runners_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"archived": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"remove_source_branch_after_merge": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			// lintignore: S031 // TODO: Resolve this tfproviderlint issue
			"push_rules": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"author_email_regex": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"branch_name_regex": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"commit_message_regex": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"commit_message_negative_regex": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_name_regex": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"commit_committer_check": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"deny_delete_tag": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"member_check": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"prevent_secrets": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"reject_unsigned_commits": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"max_file_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGitlabProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	log.Printf("[INFO] Reading Gitlab project")

	v, _ := d.GetOk("id")

	found, _, err := client.Projects.GetProject(v, nil)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", found.ID))
	d.Set("name", found.Name)
	d.Set("path", found.Path)
	d.Set("path_with_namespace", found.PathWithNamespace)
	d.Set("description", found.Description)
	d.Set("default_branch", found.DefaultBranch)
	d.Set("request_access_enabled", found.RequestAccessEnabled)
	d.Set("issues_enabled", found.IssuesEnabled)
	d.Set("merge_requests_enabled", found.MergeRequestsEnabled)
	d.Set("pipelines_enabled", found.JobsEnabled)
	d.Set("wiki_enabled", found.WikiEnabled)
	d.Set("snippets_enabled", found.SnippetsEnabled)
	d.Set("visibility_level", string(found.Visibility))
	d.Set("namespace_id", found.Namespace.ID)
	d.Set("ssh_url_to_repo", found.SSHURLToRepo)
	d.Set("http_url_to_repo", found.HTTPURLToRepo)
	d.Set("web_url", found.WebURL)
	d.Set("runners_token", found.RunnersToken)
	d.Set("archived", found.Archived)
	d.Set("remove_source_branch_after_merge", found.RemoveSourceBranchAfterMerge)

	log.Printf("[DEBUG] Reading Gitlab project %q push rules", d.Id())

	pushRules, _, err := client.Projects.GetProjectPushRules(d.Id())
	var httpError *gitlab.ErrorResponse
	if errors.As(err, &httpError) && httpError.Response.StatusCode == http.StatusNotFound {
		log.Printf("[DEBUG] Failed to get push rules for project %q: %v", d.Id(), err)
	} else if err != nil {
		return fmt.Errorf("Failed to get push rules for project %q: %w", d.Id(), err)
	}

	d.Set("push_rules", flattenProjectPushRules(pushRules)) // lintignore: XR004 // TODO: Resolve this tfproviderlint issue

	return nil
}
