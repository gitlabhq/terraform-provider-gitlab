package provider

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_branch", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_branch`" + ` resource allows to manage the lifecycle of a repository branch.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/branches.html)`,

		CreateContext: resourceGitlabBranchCreate,
		ReadContext:   resourceGitlabBranchRead,
		DeleteContext: resourceGitlabBranchDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name for this branch.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
			},
			"project": {
				Description: "The ID or full path of the project which the branch is created against.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
			},
			"ref": {
				Description: "The ref which the branch is created from.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
			},
			"web_url": {
				Description: "The url of the created branch (https).",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"default": {
				Description: "Bool, true if branch is the default branch for the project.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"can_push": {
				Description: "Bool, true if you can push to the branch.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"merged": {
				Description: "Bool, true if the branch has been merged into it's parent.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"protected": {
				Description: "Bool, true if branch has branch protection.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"developer_can_merge": {
				Description: "Bool, true if developer level access allows to merge branch.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"developer_can_push": {
				Description: "Bool, true if developer level access allows git push.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"commit": {
				Description: "The commit associated with the branch ref.",
				Type:        schema.TypeSet,
				Computed:    true,
				Set:         schema.HashResource(commitSchema),
				Elem:        commitSchema,
			},
		},
	}
})

var commitSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"id": {
			Description: "The unique id assigned to the commit by Gitlab.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"author_email": {
			Description: "The email of the author.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"author_name": {
			Description: "The name of the author.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"authored_date": {
			Description: "The date which the commit was authored (format: yyyy-MM-ddTHH:mm:ssZ).",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"committed_date": {
			Description: "The date at which the commit was pushed (format: yyyy-MM-ddTHH:mm:ssZ).",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"committer_email": {
			Description: "The email of the user that committed.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"committer_name": {
			Description: "The name of the user that committed.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"short_id": {
			Description: "The short id assigned to the commit by Gitlab.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"title": {
			Description: "The title of the commit",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"message": {
			Description: "The commit message",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"parent_ids": {
			Description: "The id of the parents of the commit",
			Type:        schema.TypeSet,
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Set:         schema.HashString,
		},
	},
}

func resourceGitlabBranchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	name := d.Get("name").(string)
	project := d.Get("project").(string)
	ref := d.Get("ref").(string)
	branchOptions := &gitlab.CreateBranchOptions{
		Branch: &name, Ref: &ref,
	}

	log.Printf("[DEBUG] create gitlab branch %s for project %s with ref %s", name, project, ref)
	branch, resp, err := client.Branches.CreateBranch(project, branchOptions, gitlab.WithContext(ctx))
	if err != nil {
		log.Printf("[DEBUG] failed to create gitlab branch %v response %v", branch, resp)
		return diag.FromErr(err)
	}
	d.Set("ref", ref)
	d.SetId(buildTwoPartID(&project, &name))
	return resourceGitlabBranchRead(ctx, d, meta)
}

func resourceGitlabBranchRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project, name, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] read gitlab branch %s", name)
	branch, resp, err := client.Branches.GetBranch(project, name, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] recieved 404 for gitlab branch %s, removing from state", name)
			d.SetId("")
			return diag.FromErr(err)
		}
		log.Printf("[DEBUG] failed to read gitlab branch %s response %v", name, resp)
		return diag.FromErr(err)
	}
	d.SetId(buildTwoPartID(&project, &name))
	d.Set("name", branch.Name)
	d.Set("project", project)
	d.Set("web_url", branch.WebURL)
	d.Set("default", branch.Default)
	d.Set("can_push", branch.CanPush)
	d.Set("merged", branch.Merged)
	d.Set("developer_can_merge", branch.DevelopersCanMerge)
	d.Set("developer_can_push", branch.DevelopersCanPush)
	d.Set("protected", branch.Protected)
	commit := flattenCommit(branch.Commit)
	if err := d.Set("commit", commit); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceGitlabBranchDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project, name, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] delete gitlab branch %s", name)
	resp, err := client.Branches.DeleteBranch(project, name, gitlab.WithContext(ctx))
	if err != nil {
		log.Printf("[DEBUG] failed to delete gitlab branch %s response %v", name, resp)
		return diag.FromErr(err)
	}
	return nil
}

func flattenCommit(commit *gitlab.Commit) (values []map[string]interface{}) {
	if commit == nil {
		return []map[string]interface{}{}
	}
	return []map[string]interface{}{
		{
			"id":              commit.ID,
			"short_id":        commit.ShortID,
			"title":           commit.Title,
			"author_name":     commit.AuthorName,
			"author_email":    commit.AuthorEmail,
			"authored_date":   commit.AuthoredDate.Format(time.RFC3339),
			"committed_date":  commit.CommittedDate.Format(time.RFC3339),
			"committer_email": commit.CommitterEmail,
			"committer_name":  commit.CommitterName,
			"message":         commit.Message,
			"parent_ids":      commit.ParentIDs,
		},
	}
}
