package gitlab

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabProjectBadge() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabProjectBadgeCreate,
		Read:   resourceGitlabProjectBadgeRead,
		Update: resourceGitlabProjectBadgeUpdate,
		Delete: resourceGitlabProjectBadgeDelete,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"link_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"image_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rendered_link_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rendered_image_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGitlabProjectBadgeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	options := &gitlab.AddProjectBadgeOptions{
		LinkURL:  gitlab.String(d.Get("link_url").(string)),
		ImageURL: gitlab.String(d.Get("image_url").(string)),
	}

	log.Printf("[DEBUG] create gitlab project badge %q / %q", *options.LinkURL, *options.ImageURL)

	badge, _, err := client.ProjectBadges.AddProjectBadge(project, options)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", badge.ID))

	return resourceGitlabProjectBadgeRead(d, meta)
}

func resourceGitlabProjectBadgeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	badgeId, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] read gitlab project badge %s/%d", project, badgeId)

	badge, _, err := client.ProjectBadges.GetProjectBadge(project, badgeId)
	if err != nil {
		return err
	}

	d.Set("link_url", badge.LinkURL)
	d.Set("image_url", badge.ImageURL)
	d.Set("rendered_link_url", badge.RenderedLinkURL)
	d.Set("rendered_image_url", badge.RenderedImageURL)
	return nil
}

func resourceGitlabProjectBadgeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	badgeId, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	options := &gitlab.EditProjectBadgeOptions{
		LinkURL:  gitlab.String(d.Get("link_url").(string)),
		ImageURL: gitlab.String(d.Get("image_url").(string)),
	}

	log.Printf("[DEBUG] update gitlab project badge %s", d.Id())

	_, _, err = client.ProjectBadges.EditProjectBadge(project, badgeId, options)
	if err != nil {
		return err
	}

	return resourceGitlabProjectBadgeRead(d, meta)
}

func resourceGitlabProjectBadgeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	badgeId, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Delete gitlab project badge %s", d.Id())

	_, err = client.ProjectBadges.DeleteProjectBadge(project, badgeId)
	return err
}
