package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_project_issue", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_project_issue`" + ` data source allows to retrieve details about an issue in a project.

**Upstream API**: [GitLab API docs](https://docs.gitlab.com/ee/api/issues.html)`,

		ReadContext: dataSourceGitlabProjectIssueRead,
		Schema:      datasourceSchemaFromResourceSchema(gitlabProjectIssueGetSchema(), []string{"project", "iid"}, nil),
	}
})

func dataSourceGitlabProjectIssueRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	issueIID := d.Get("iid").(int)

	issue, _, err := client.Issues.GetIssue(project, issueIID, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resourceGitLabProjectIssueBuildId(project, issueIID))
	stateMap := gitlabProjectIssueToStateMap(project, issue)

	if err := setStateMapInResourceData(stateMap, d); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
