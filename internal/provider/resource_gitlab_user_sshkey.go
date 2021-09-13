package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_user_sshkey", func() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGitlabUserSSHKeyCreate,
		ReadContext:   resourceGitlabUserSSHKeyRead,
		UpdateContext: resourceGitlabUserSSHKeyUpdate,
		DeleteContext: resourceGitlabUserSSHKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGitlabUserSSHKeyImporter,
		},

		Schema: map[string]*schema.Schema{
			"title": {
				Description: "The title of the ssh key.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"key": {
				Description: "The ssh key.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"created_at": {
				Description: "Create time.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"user_id": {
				Description: "The id of the user to add the ssh key to.",
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
			},
		},
	}
})

func resourceGitlabUserSSHKeyImporter(ctx context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	s := strings.Split(d.Id(), ":")
	if len(s) != 2 {
		d.SetId("")
		return nil, fmt.Errorf("Invalid SSH Key import format; expected '{user_id}:{key_id}'")
	}

	userID, err := strconv.Atoi(s[0])
	if err != nil {
		return nil, fmt.Errorf("Invalid SSH Key import format; expected '{user_id}:{key_id}'")
	}

	d.Set("user_id", userID)
	d.SetId(s[1])

	return []*schema.ResourceData{d}, nil
}

func resourceGitlabUserSSHKeySetToState(d *schema.ResourceData, key *gitlab.SSHKey) {
	d.Set("title", key.Title)
	d.Set("key", key.Key)
	d.Set("created_at", key.CreatedAt.String())
}

func resourceGitlabUserSSHKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	userID := d.Get("user_id").(int)

	options := &gitlab.AddSSHKeyOptions{
		Title: gitlab.String(d.Get("title").(string)),
		Key:   gitlab.String(d.Get("key").(string)),
	}

	key, _, err := client.Users.AddSSHKeyForUser(userID, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", key.ID))

	return resourceGitlabUserSSHKeyRead(ctx, d, meta)
}

func resourceGitlabUserSSHKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	id, _ := strconv.Atoi(d.Id())
	userID := d.Get("user_id").(int)

	keys, _, err := client.Users.ListSSHKeysForUser(userID, &gitlab.ListSSHKeysForUserOptions{}, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	var key *gitlab.SSHKey

	for _, k := range keys {
		if k.ID == id {
			key = k
			break
		}
	}

	if key == nil {
		return diag.Errorf("Could not find sshkey %d for user %d", id, userID)
	}

	resourceGitlabUserSSHKeySetToState(d, key)
	return nil
}

func resourceGitlabUserSSHKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d := resourceGitlabUserSSHKeyDelete(ctx, d, meta); d.HasError() {
		return d
	}
	if d := resourceGitlabUserSSHKeyCreate(ctx, d, meta); d.HasError() {
		return d
	}
	return resourceGitlabUserSSHKeyRead(ctx, d, meta)
}

func resourceGitlabUserSSHKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	id, _ := strconv.Atoi(d.Id())
	userID := d.Get("user_id").(int)

	if _, err := client.Users.DeleteSSHKeyForUser(userID, id, gitlab.WithContext(ctx)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
