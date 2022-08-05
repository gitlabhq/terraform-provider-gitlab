package provider

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_user_gpgkey", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`" + `gitlab_user_gpgkey` + "`" + ` resource allows to manage the lifecycle of a GPG key assigned to the current user or a specific user.
		
-> Managing GPG keys for arbitrary users requires admin privileges.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/users.html#get-a-specific-gpg-key)`,

		CreateContext: resourceGitlabUserGPGKeyCreate,
		ReadContext:   resourceGitlabUserGPGKeyRead,
		DeleteContext: resourceGitlabUserGPGKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"user_id": {
				Description: "The ID of the user to add the GPG key to. If this field is omitted, this resource manages a GPG key for the current user. Otherwise, this resource manages a GPG key for the speicifed user, and an admin token is required.",
				Type:        schema.TypeInt,
				ForceNew:    true,
				Optional:    true,
			},
			"key": {
				Description: "The armored GPG public key.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					strippedOld := strings.TrimSpace(old)
					strippedNew := strings.TrimSpace(new)
					return strippedOld == strippedNew
				},
			},
			"key_id": {
				Description: "The ID of the GPG key.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"created_at": {
				Description: "The time when this key was created in GitLab.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
})

func resourceGitlabUserGPGKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	options := &gitlab.AddGPGKeyOptions{
		Key: gitlab.String(strings.TrimSpace(d.Get("key").(string))),
	}

	var isAdmin bool
	var key *gitlab.GPGKey
	var err error
	userID, userIDOk := d.GetOk("user_id")
	if userIDOk {
		isAdmin, err = isCurrentUserAdmin(ctx, client)
		if err != nil {
			return diag.Errorf("failed to check if user is admin for configuring GPG keys for a user")
		}
		if !isAdmin {
			return diag.Errorf("current user needs to be admin for configuring GPG keys for a user")
		}
		key, _, err = client.Users.AddGPGKeyForUser(userID.(int), options, gitlab.WithContext(ctx))
	} else {
		key, _, err = client.Users.AddGPGKey(options, gitlab.WithContext(ctx))
	}
	if err != nil {
		return diag.FromErr(err)
	}

	keyIDForID := fmt.Sprintf("%d", key.ID)
	if userIDOk {
		userIDForID := fmt.Sprintf("%d", userID)
		d.SetId(buildTwoPartID(&userIDForID, &keyIDForID))
	} else {
		d.SetId(keyIDForID)
	}
	return resourceGitlabUserGPGKeyRead(ctx, d, meta)
}

func resourceGitlabUserGPGKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	userID, keyID, err := resourceGitlabUserGPGKeyParseID(d.Id())
	if err != nil {
		return diag.Errorf("unable to parse user GPG key resource id: %s: %v", d.Id(), err)
	}

	var key *gitlab.GPGKey
	if userID != 0 {
		key, _, err = client.Users.GetGPGKeyForUser(userID, keyID, gitlab.WithContext(ctx))
	} else {
		key, _, err = client.Users.GetGPGKey(keyID, gitlab.WithContext(ctx))
	}
	if err != nil {
		if is404(err) {
			log.Printf("Could not find GPG key %d for user %d, removing from state", keyID, userID)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if userID != 0 {
		d.Set("user_id", userID)
	}
	d.Set("key_id", keyID)
	d.Set("key", strings.TrimSpace(key.Key))
	d.Set("created_at", key.CreatedAt.Format(time.RFC3339))
	return nil
}

func resourceGitlabUserGPGKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	var isAdmin bool
	_, keyID, err := resourceGitlabUserGPGKeyParseID(d.Id())
	if err != nil {
		return diag.Errorf("unable to parse user GPG key resource id: %s: %v", d.Id(), err)
	}

	if userID, ok := d.GetOk("user_id"); ok {
		isAdmin, err = isCurrentUserAdmin(ctx, client)
		if err != nil {
			return diag.Errorf("failed to check if user is admin for configuring GPG keys for a user")
		}
		if !isAdmin {
			return diag.Errorf("current user needs to be admin for configuring GPG keys for a user")
		}
		_, err = client.Users.DeleteGPGKeyForUser(userID.(int), keyID, gitlab.WithContext(ctx))
	} else {
		_, err = client.Users.DeleteGPGKey(keyID, gitlab.WithContext(ctx))
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGitlabUserGPGKeyParseID(id string) (int, int, error) {
	userIDFromID, keyIDFromID, err := parseTwoPartID(id)
	if err != nil {
		keyID, errKeyID := strconv.Atoi(id)
		if errKeyID != nil {
			return 0, 0, err
		} else {
			return 0, keyID, nil
		}
	}
	userID, err := strconv.Atoi(userIDFromID)
	if err != nil {
		return 0, 0, err
	}
	keyID, err := strconv.Atoi(keyIDFromID)
	if err != nil {
		return 0, 0, err
	}

	return userID, keyID, nil
}
