package gitlab

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabGroupBadge() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGitlabGroupBadgeCreate,
		ReadContext:   resourceGitlabGroupBadgeRead,
		UpdateContext: resourceGitlabGroupBadgeUpdate,
		DeleteContext: resourceGitlabGroupBadgeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"group": {
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

func resourceGitlabGroupBadgeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	groupID := d.Get("group").(string)
	options := &gitlab.AddGroupBadgeOptions{
		LinkURL:  gitlab.String(d.Get("link_url").(string)),
		ImageURL: gitlab.String(d.Get("image_url").(string)),
	}

	log.Printf("[DEBUG] create gitlab group variable %s/%s", *options.LinkURL, *options.ImageURL)

	badge, _, err := client.GroupBadges.AddGroupBadge(groupID, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	badgeID := strconv.Itoa(badge.ID)

	d.SetId(buildTwoPartID(&groupID, &badgeID))

	return resourceGitlabGroupBadgeRead(ctx, d, meta)
}

func resourceGitlabGroupBadgeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	ids := strings.Split(d.Id(), ":")
	groupID := ids[0]
	badgeID, err := strconv.Atoi(ids[1])
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] read gitlab group badge %s/%d", groupID, badgeID)

	badge, _, err := client.GroupBadges.GetGroupBadge(groupID, badgeID, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	resourceGitlabGroupBadgeSetToState(d, badge, &groupID)
	return nil
}

func resourceGitlabGroupBadgeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	ids := strings.Split(d.Id(), ":")
	groupID := ids[0]
	badgeID, err := strconv.Atoi(ids[1])
	if err != nil {
		return diag.FromErr(err)
	}

	options := &gitlab.EditGroupBadgeOptions{
		LinkURL:  gitlab.String(d.Get("link_url").(string)),
		ImageURL: gitlab.String(d.Get("image_url").(string)),
	}

	log.Printf("[DEBUG] update gitlab group badge %s/%d", groupID, badgeID)

	_, _, err = client.GroupBadges.EditGroupBadge(groupID, badgeID, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGitlabGroupBadgeRead(ctx, d, meta)
}

func resourceGitlabGroupBadgeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	ids := strings.Split(d.Id(), ":")
	groupID := ids[0]
	badgeID, err := strconv.Atoi(ids[1])
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Delete gitlab group badge %s/%d", groupID, badgeID)

	_, err = client.GroupBadges.DeleteGroupBadge(groupID, badgeID, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGitlabGroupBadgeSetToState(d *schema.ResourceData, badge *gitlab.GroupBadge, groupID *string) {
	d.Set("link_url", badge.LinkURL)
	d.Set("image_url", badge.ImageURL)
	d.Set("rendered_link_url", badge.RenderedLinkURL)
	d.Set("rendered_image_url", badge.RenderedImageURL)
	d.Set("group", groupID)
}
