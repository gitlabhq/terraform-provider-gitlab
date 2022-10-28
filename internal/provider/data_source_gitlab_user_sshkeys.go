package provider

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_user_sshkeys", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_user_sshkeys`" + ` data source allows a list of SSH keys to be retrieved by either the user ID or username.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/users.html#list-ssh-keys-for-user)`,

		ReadContext: dataSourceGitlabUserKeysRead,
		Schema: map[string]*schema.Schema{
			"user_id": {
				Description: "ID of the user to get the SSH keys for.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ConflictsWith: []string{
					"username",
				},
			},
			"username": {
				Description: "Username of the user to get the SSH keys for.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ConflictsWith: []string{
					"user_id",
				},
			},
			"keys": {
				Description: "The user's keys.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: datasourceSchemaFromResourceSchema(gitlabUserSSHKeySchema(), nil, nil),
				},
			},
		},
	}
})

func dataSourceGitlabUserKeysRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	var keys []*gitlab.SSHKey
	var err error

	log.Printf("[INFO] Reading Gitlab user")

	userIDData, userIDOk := d.GetOk("user_id")
	usernameData, usernameOk := d.GetOk("username")

	if userIDOk {
		// Get SSH keys by user ID
		keys, _, err = client.Users.ListSSHKeysForUser(userIDData.(int), &gitlab.ListSSHKeysForUserOptions{}, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	} else if usernameOk {
		username := strings.ToLower(usernameData.(string))
		// Get SSH keys by username
		keys, _, err = client.Users.ListSSHKeysForUser(username, &gitlab.ListSSHKeysForUserOptions{}, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.Errorf("one and only one of user_id or username must be set")
	}

	d.SetId(fmt.Sprintf("%d", userIDData))
	if err := d.Set("keys", flattenSSHKeysForState(keys)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func flattenSSHKeysForState(keys []*gitlab.SSHKey) (values []map[string]interface{}) {
	for _, key := range keys {
		values = append(values, gitlabUserKeyToStateMap(key))
	}
	return values
}
