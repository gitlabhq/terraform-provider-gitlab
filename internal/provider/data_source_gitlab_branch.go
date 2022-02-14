package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_branch", func() *schema.Resource {
	return &schema.Resource{
		Description: "Provide details about a gitlab project branch \n\n" +
			"> **Note**: name and project must be specified.",

		ReadContext: dataSourceGitlabBranchRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the branch.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"project": {
				Description: "The full path or id of the group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"web_url": {
				Description: "The url of the created branch (https)",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"default": {
				Description: "Bool, true if branch is the default branch for the project",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"can_push": {
				Description: "Bool, true if you can push to the branch",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"protected": {
				Description: "Bool, true if branch has branch protection",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"merged": {
				Description: "Bool, true if the branch has been merged into it's parent",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"developer_can_merge": {
				Description: "Bool, true if developer level access allows to merge branch",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"developer_can_push": {
				Description: "Bool, true if developer level access allows git push",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"commit": {
				Description: "The commit associated with the branch ref",
				Type:        schema.TypeSet,
				Computed:    true,
				Set:         schema.HashResource(commitSchema),
				Elem:        commitSchema,
			},
		},
	}
})

func dataSourceGitlabBranchRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	name := d.Get("name").(string)
	project := d.Get("project").(string)
	log.Printf("[DEBUG] read gitlab branch %s", name)
	branch, resp, err := client.Branches.GetBranch(project, name, gitlab.WithContext(ctx))
	if err != nil {
		log.Printf("[DEBUG] failed to read gitlab branch %s response %v", name, resp)
		return diag.FromErr(err)
	}

	d.SetId(buildTwoPartID(&project, &name))
	d.Set("name", branch.Name)
	d.Set("project", project)
	d.Set("web_url", branch.WebURL)
	d.Set("default", branch.Default)
	d.Set("can_push", branch.CanPush)
	d.Set("protected", branch.Protected)
	d.Set("merged", branch.Merged)
	d.Set("developer_can_merge", branch.DevelopersCanMerge)
	d.Set("developer_can_push", branch.DevelopersCanPush)
	commit := flattenCommit(branch.Commit)
	if err := d.Set("commit", commit); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
