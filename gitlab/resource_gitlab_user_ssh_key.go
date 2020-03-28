package gitlab

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabUserSSHKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabUserSSHKeyCreate,
		Read:   resourceGitlabUserSSHKeyRead,
		Delete: resourceGitlabUserSSHKeyDelete,

		Schema: map[string]*schema.Schema{
			"title": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, oldV, newV string, d *schema.ResourceData) bool {
					newTrimmed := strings.TrimSpace(newV)
					return oldV == newTrimmed
				},
			},
			"key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGitlabUserSSHKeySetToState(d *schema.ResourceData, userkey *gitlab.SSHKey) {
	d.Set("key_id", userkey.ID)
	d.Set("title", userkey.Title)
	d.Set("key", userkey.Key)
	d.Set("created_at", userkey.CreatedAt)
}

func resourceGitlabUserSSHKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	options := &gitlab.AddSSHKeyOptions{
		Title: gitlab.String(d.Get("title").(string)),
		Key:   gitlab.String(d.Get("key").(string)),
	}

	log.Printf("[DEBUG] create gitlab user SSH key %s", *options.Title)
	userKey, _, err := client.Users.AddSSHKey(options)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", userKey.ID))
	d.Set("title", userKey.Title)
	return resourceGitlabUserSSHKeyRead(d, meta)
}

func resourceGitlabUserSSHKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	log.Printf("[DEBUG] read gitlab user SSH Key %s", d.Id())
	id, _ := strconv.Atoi(d.Id())

	key, _, err := client.Users.GetSSHKey(id)
	if err != nil {
		return err
	}
	resourceGitlabUserSSHKeySetToState(d, key)
	return nil
}

func resourceGitlabUserSSHKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	log.Printf("[DEBUG] destroy gitlab user SSH Key %s", d.Id())

	id, _ := strconv.Atoi(d.Id())

	_, err := client.Users.DeleteSSHKey(id)
	log.Printf("[DEBUG] Delete gitlab user SSH Key%s", err)
	return nil
}
