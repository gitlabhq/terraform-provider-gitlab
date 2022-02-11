package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_branch", func() *schema.Resource {
	return &schema.Resource{
		Description: "This resource allows you to create and manage GitLab branch.",

		CreateContext: resourceGitlabBranchCreate,
		ReadContext:   resourceGitlabBranchRead,
		DeleteContext: resourceGitlabBranchDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Description: "The ID of the project",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"branch": {
				Description: "The name of the new branch",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"ref": {
				Description: "The branch from which the new branch is created",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
})

func resourceGitlabBranchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	options := &gitlab.CreateBranchOptions{
		Branch: gitlab.String(d.Get("branch").(string)),
		Ref:    gitlab.String(d.Get("ref").(string)),
	}

	log.Printf("[DEBUG] create gitlab branch %q", *options.Branch)

	branch, _, err := client.Branches.CreateBranch(project, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(buildTwoPartID(&project, &branch.Name))

	return resourceGitlabBranchRead(ctx, d, meta)
}

func resourceGitlabBranchRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project, branchName, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] read gitlab branch %s", branchName)

	branch, _, err := client.Branches.GetBranch(project, branchName, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] gitlab branch not found %s/%s", project, branchName)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	c := &gitlab.GetCommitRefsOptions{
		Type: gitlab.String("branch"),
	}
	commitRefs, _, err := client.Commits.GetCommitRefs(project, branch.Commit.ID, c, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}
	for _, br := range commitRefs {
		if br.Name != branch.Name {
			d.Set("ref", br.Name)
		}
	}

	d.Set("project", project)
	d.Set("branch", branch.Name)

	d.SetId(buildTwoPartID(&project, &branch.Name))

	return nil
}

func resourceGitlabBranchDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project, branchName, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Delete gitlab branch %s", d.Id())

	resp, err := client.Branches.DeleteBranch(project, branchName, gitlab.WithContext(ctx))
	if err != nil {
		return diag.Errorf("%s failed to delete branch: (%s) %v", branchName, resp.Status, err)
	}

	return nil
}
