package provider

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/hashstructure"
	"github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_project_tags", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_project_tags`" + ` data source allows details of project tags to be retrieved by some search criteria.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/tags.html#list-project-repository-tags)`,

		ReadContext: dataSourceGitlabProjectTagsRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Description: "The ID or URL-encoded path of the project owned by the authenticated user.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"order_by": {
				Description: "Return tags ordered by `name` or `updated` fields. Default is `updated`.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"sort": {
				Description: "Return tags sorted in `asc` or `desc` order. Default is `desc`.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"search": {
				Description: "Return list of tags matching the search criteria. You can use `^term` and `term$` to find tags that begin and end with `term` respectively. No other regular expressions are supported.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"tags": {
				Description: "List of repository tags from a project.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: datasourceSchemaFromResourceSchema(gitlabProjectTagGetSchema(), nil, nil),
				},
			},
		},
	}
})

func dataSourceGitlabProjectTagsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	project := d.Get("project").(string)
	options := gitlab.ListTagsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 20,
			Page:    1,
		},
	}

	if v, ok := d.GetOk("order_by"); ok {
		options.OrderBy = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("sort"); ok {
		options.Sort = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("search"); ok {
		options.Search = gitlab.String(v.(string))
	}

	optionsHash, err := hashstructure.Hash(&options, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var tags []*gitlab.Tag
	for options.Page != 0 {
		paginatedTags, resp, err := client.Tags.ListTags(project, &options, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}

		tags = append(tags, paginatedTags...)
		options.Page = resp.NextPage
	}

	log.Printf("[DEBUG] get gitlab tags from project: %s", project)
	d.SetId(fmt.Sprintf("%s:%d", project, optionsHash))
	d.Set("project", project)
	d.Set("order_by", options.OrderBy)
	d.Set("sort", options.Sort)
	d.Set("search", options.Search)
	if err = d.Set("tags", flattenDataTags(tags)); err != nil {
		return diag.Errorf("Failed to set tags to state: %v", err)
	}
	return nil
}

func flattenDataTags(tags []*gitlab.Tag) (values []map[string]interface{}) {
	for _, tag := range tags {
		values = append(values, map[string]interface{}{
			"commit":    flattenCommit(tag.Commit),
			"release":   flattenReleaseNote(tag.Release),
			"name":      tag.Name,
			"target":    tag.Target,
			"message":   tag.Message,
			"protected": tag.Protected,
		})
	}
	return values
}
