package gitlab

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabDeployKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabDeployKeyCreate,
		Read:   resourceGitlabDeployKeyRead,
		Delete: resourceGitlabDeployKeyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGitlabDeployKeyStateImporter,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"title": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					var oldPubKey, newPubKey ssh.PublicKey
					var err error

					switch new == "" {
					case true:
						return false
					case false:
						newPubKey, _, _, _, err = ssh.ParseAuthorizedKey([]byte(new))
						if err != nil {
							panic(err)
						}
					}

					switch old == "" {
					case true:
						return false
					case false:
						oldPubKey, _, _, _, err = ssh.ParseAuthorizedKey([]byte(old))
						if err != nil {
							panic(err)
						}
					}
					return string(bytes.TrimSpace(oldPubKey.Marshal())) == string(bytes.TrimSpace(newPubKey.Marshal()))
				},
			},
			"can_push": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
		},
	}
}

func resourceGitlabDeployKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	options := &gitlab.AddDeployKeyOptions{
		Title:   gitlab.String(d.Get("title").(string)),
		Key:     gitlab.String(strings.TrimSpace(d.Get("key").(string))),
		CanPush: gitlab.Bool(d.Get("can_push").(bool)),
	}

	log.Printf("[DEBUG] create gitlab deployment key %s", *options.Title)

	deployKey, _, err := client.DeployKeys.AddDeployKey(project, options)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", deployKey.ID))

	return resourceGitlabDeployKeyRead(d, meta)
}

func resourceGitlabDeployKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	deployKeyID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] read gitlab deploy key %s/%d", project, deployKeyID)

	deployKey, _, err := client.DeployKeys.GetDeployKey(project, deployKeyID)
	if err != nil {
		return err
	}

	d.Set("title", deployKey.Title)
	d.Set("key", deployKey.Key)
	d.Set("can_push", deployKey.CanPush)
	return nil
}

func resourceGitlabDeployKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	deployKeyID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Delete gitlab deploy key %s", d.Id())

	response, err := client.DeployKeys.DeleteDeployKey(project, deployKeyID)

	// HTTP 204 is success with no body
	if response.StatusCode == 204 {
		return nil
	}
	return err
}

func resourceGitlabDeployKeyStateImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := strings.Split(d.Id(), ":")
	if len(s) != 2 {
		d.SetId("")
		return nil, fmt.Errorf("Invalid Deploy Key import format; expected '{project_id}:{deploy_key_id}'")
	}
	project, id := s[0], s[1]

	d.SetId(id)
	d.Set("project", project)

	return []*schema.ResourceData{d}, nil
}
