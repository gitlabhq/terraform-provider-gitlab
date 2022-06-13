package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_release_link", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_release_link`" + ` data source allows get details of a release link.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/releases/links.html)`,

		ReadContext: dataSourceGitlabReleaseLinkRead,
		Schema:      datasourceSchemaFromResourceSchema(gitlabReleaseLinkGetSchema(), []string{"project", "tag_name", "link_id"}, nil),
	}
})

func dataSourceGitlabReleaseLinkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	tagName := d.Get("tag_name").(string)
	linkID := d.Get("link_id").(int)

	releaseLink, _, err := client.ReleaseLinks.GetReleaseLink(project, tagName, linkID, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resourceGitLabReleaseLinkBuildId(project, tagName, linkID))
	stateMap := gitlabReleaseLinkToStateMap(project, tagName, releaseLink)

	if err := setStateMapInResourceData(stateMap, d); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
