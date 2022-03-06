package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_project_tag", func() *schema.Resource {
	return &schema.Resource{
		Description: `This resource allows you to create and manage GitLab tags.

**Upstream API** : [GitLab API docs](https://docs.gitlab.com/ee/api/tags.html)`,

		CreateContext: resourceGitlabTagCreate,
		ReadContext:   resourceGitlabTagRead,
		DeleteContext: resourceGitlabTagDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of a tag.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
			},
			"project": {
				Description: "The ID or URL-encoded path of the project owned by the authenticated user.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
			},
			"ref": {
				Description: "Create tag using commit SHA, another tag name, or branch name. This attribute is not available for imported resources.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
			},
			"message": {
				Description: "Creates annotated tag.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
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
				Description: "The commit associated with the tag.",
				Type:        schema.TypeSet,
				Computed:    true,
				Set:         schema.HashResource(commitSchema),
				Elem:        commitSchema,
			},
		},
	}
})

var releaseNoteSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"tag_name": {
			Description: "The name of the tag.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"description": {
			Description: "The description of release.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	},
}

func resourceGitlabTagCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	name := d.Get("name").(string)
	project := d.Get("project").(string)
	ref := d.Get("ref").(string)
	message := d.Get("message").(string)
	tagOptions := &gitlab.CreateTagOptions{
		TagName: &name, Ref: &ref, Message: &message,
	}

	log.Printf("[DEBUG] create gitlab tag %s/%s with ref %s", project, name, ref)
	_, resp, err := client.Tags.CreateTag(project, tagOptions, gitlab.WithContext(ctx))
	if err != nil {
		log.Printf("[DEBUG] failed to create gitlab tag %s/%s response %v", project, name, resp)
		return diag.FromErr(err)
	}
	d.SetId(buildTwoPartID(&project, &name))
	d.Set("ref", ref)
	return resourceGitlabTagRead(ctx, d, meta)
}

func resourceGitlabTagRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project, name, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] read gitlab tag %s/%s", project, name)
	tag, resp, err := client.Tags.GetTag(project, name, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] recieved 404 for gitlab tag %s/%s, removing from state", project, name)
			d.SetId("")
			return diag.FromErr(err)
		}
		log.Printf("[DEBUG] failed to read gitlab tag %s/%s response %v", project, name, resp)
		return diag.FromErr(err)
	}
	d.Set("name", tag.Name)
	d.Set("project", project)
	d.Set("message", tag.Message)
	d.Set("protected", tag.Protected)
	d.Set("target", tag.Target)
	releaseNote := flattenReleaseNote(tag.Release)
	if err := d.Set("release", releaseNote); err != nil {
		return diag.FromErr(err)
	}
	commit := flattenCommit(tag.Commit)
	if err := d.Set("commit", commit); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceGitlabTagDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project, name, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] delete gitlab tag %s/%s", project, name)
	resp, err := client.Tags.DeleteTag(project, name, gitlab.WithContext(ctx))
	if err != nil {
		log.Printf("[DEBUG] failed to delete gitlab tag %s/%s response %v", project, name, resp)
		return diag.FromErr(err)
	}
	return nil
}

func flattenReleaseNote(releaseNote *gitlab.ReleaseNote) (values []map[string]interface{}) {
	if releaseNote == nil {
		return []map[string]interface{}{}
	}
	return []map[string]interface{}{
		{
			"tag_name":    releaseNote.TagName,
			"description": releaseNote.Description,
		},
	}
}
