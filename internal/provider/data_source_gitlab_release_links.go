package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_release_links", func() *schema.Resource {

	return &schema.Resource{
		Description: `The ` + "`gitlab_release_links`" + ` data source allows get details of release links.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/releases/links.html)`,

		ReadContext: dataSourceGitlabReleaseLinksRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Description: "The ID or full path to the project.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"tag_name": {
				Description: "The tag associated with the Release.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"release_links": {
				Description: "List of release links",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: datasourceSchemaFromResourceSchema(gitlabReleaseLinkGetSchema(), nil, nil),
				},
			},
		},
	}
})

func dataSourceGitlabReleaseLinksRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	project := d.Get("project").(string)
	tagName := d.Get("tag_name").(string)
	options := gitlab.ListReleaseLinksOptions(
		gitlab.ListOptions{
			PerPage: 20,
			Page:    1,
		})

	var releaseLinks []*gitlab.ReleaseLink
	for options.Page != 0 {
		paginatedReleaseLinks, resp, err := client.ReleaseLinks.ListReleaseLinks(project, tagName, &options, gitlab.WithContext(ctx))
		if err != nil {
			if is404(err) && (options.Page == 1) {
				break
			} else {
				return diag.FromErr(err)
			}
		}
		releaseLinks = append(releaseLinks, paginatedReleaseLinks...)
		options.Page = resp.NextPage
	}

	log.Printf("[DEBUG] get list release links project/tagName: %s/%s", project, tagName)
	d.SetId(buildTwoPartID(&project, &tagName))
	if err := d.Set("release_links", flattenGitlabReleaseLinks(project, tagName, releaseLinks)); err != nil {
		return diag.Errorf("Failed to set release links to state: %v", err)
	}

	return nil
}

func flattenGitlabReleaseLinks(project string, tagName string, releaseLinks []*gitlab.ReleaseLink) (values []map[string]interface{}) {
	for _, releaseLink := range releaseLinks {
		values = append(values, gitlabReleaseLinkToStateMap(project, tagName, releaseLink))
	}
	return values
}
