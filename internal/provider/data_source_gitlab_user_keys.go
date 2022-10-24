package provider

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_user_keys", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_user_keys`" + ` data source allows a list of SSH keys to be retrieved by either the user ID or username.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/users.html#list-ssh-keys-for-user)`,

		ReadContext: dataSourceGitlabUserKeysRead,
		Schema:      gitlabUserKeysSchema(),
	}
})

func dataSourceGitlabUserKeysRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	var keys []*gitlab.SSHKey
	var err error

	log.Printf("[INFO] Reading Gitlab user")

	userIDData, userIDOk := d.GetOk("id_or_username")

	if userIDOk {
		// Get SSH keys by user
		keys, _, err = client.Users.ListSSHKeysForUser(userIDData, &gitlab.ListSSHKeysForUserOptions{}, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.Errorf("one and only one of user_id or username must be set")
	}

	d.SetId(fmt.Sprintf("%d", userIDData))
	d.Set("keys", flattenSSHKeysForState(keys))

	return nil
}

func flattenSSHKeysForState(keys []*gitlab.SSHKey) (values []map[string]interface{}) {
	for _, key := range keys {
		values = append(values, gitlabUserKeyToStateMap(key))
	}
	return values
}
