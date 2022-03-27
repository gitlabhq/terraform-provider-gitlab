package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_project_tag", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_project_tag`" + ` data source allows details of a project tag to be retrieved by its name.

**Upstream API**: [GitLab API docs](https://docs.gitlab.com/ee/api/tags.html)`,

		ReadContext: dataSourceGitlabProjectTagRead,
		Schema: constructSchema(
			datasourceSchemaFromResourceSchema(gitlabProjectTagGetSchema(), []string{"name"}, nil),
			map[string]*schema.Schema{
				"project": {
					Description: "The ID or URL-encoded path of the project owned by the authenticated user.",
					Type:        schema.TypeString,
					Required:    true,
				},
			},
		),
	}
})

func dataSourceGitlabProjectTagRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	name := d.Get("name").(string)
	project := d.Get("project").(string)
	log.Printf("[DEBUG] read gitlab tag %s/%s", project, name)
	tag, resp, err := client.Tags.GetTag(project, name, gitlab.WithContext(ctx))
	if err != nil {
		log.Printf("[DEBUG] failed to read gitlab tag %s/%s response %v", project, name, resp)
		return diag.FromErr(err)
	}

	d.SetId(buildTwoPartID(&project, &name))
	d.Set("name", tag.Name)
	d.Set("project", project)
	d.Set("message", tag.Message)
	d.Set("protected", tag.Protected)
	d.Set("target", tag.Target)
	releaseNote := flattenReleaseNote(tag.Release)
	if err := d.Set("release", releaseNote); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("commit", flattenCommit(tag.Commit)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
