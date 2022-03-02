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
		Description: `Provide details about a gitlab project tag
				
**Upstream API** : [GitLab API docs](https://docs.gitlab.com/ee/api/tags.html)`,

		ReadContext: dataSourceGitlabTagRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of a tag.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"project": {
				Description: "The ID or URL-encoded path of the project owned by the authenticated user.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"message": {
				Description: "Creates annotated tag.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"protected": {
				Description: "Bool, true if tag has tag protection.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"target": {
				Description: "The unique id assigned to the commit by Gitlab.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"release": {
				Description: "The release associated with the tag.",
				Type:        schema.TypeSet,
				Computed:    true,
				Set:         schema.HashResource(releaseNoteSchema),
				Elem:        releaseNoteSchema,
			},
			"commit": {
				Description: "The commit associated with the tag ref.",
				Type:        schema.TypeSet,
				Computed:    true,
				Set:         schema.HashResource(commitSchema),
				Elem:        commitSchema,
			},
		},
	}
})

func dataSourceGitlabTagRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
