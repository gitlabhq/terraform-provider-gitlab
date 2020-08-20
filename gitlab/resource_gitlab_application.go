package gitlab

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

// https://docs.gitlab.com/ce/api/applications.html
func resourceGitlabApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabApplicationCreate,
		Read:   resourceGitlabApplicationRead,
		Update: resourceGitlabApplicationUpdate,
		Delete: resourceGitlabApplicationDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scopes": {
				Type:     schema.TypeSet,
				ForceNew: true,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"redirect_uri": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"confidential": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
			"application_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"callback_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGitlabApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	uris := strings.Join(*stringSetToStringSlice(d.Get("redirect_uri").(*schema.Set)), " ")
	scopes := strings.Join(*stringSetToStringSlice(d.Get("scopes").(*schema.Set)), " ")

	options := &gitlab.CreateApplicationOptions{
		Name:        gitlab.String(d.Get("name").(string)),
		RedirectURI: gitlab.String(uris),
		Scopes:      gitlab.String(scopes),
	}

	log.Printf("[DEBUG] create application %s", *options.Name)

	app, _, err := client.Applications.CreateApplication(options)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(app.ID))
	d.Set("application_id", app.ApplicationID)
	d.Set("secret", app.Secret)

	return resourceGitlabApplicationRead(d, meta)
}

func resourceGitlabApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	appID := d.Get("application_id").(string)

	log.Printf("[DEBUG] read gitlab application by id %s", appID)

	page := 1
	appsLen := 0
	for page == 1 || appsLen != 0 {
		apps, _, err := client.Applications.ListApplications(&gitlab.ListApplicationsOptions{Page: page})
		if err != nil {
			return err
		}
		for _, app := range apps {
			if app.ApplicationID == appID {
				d.Set("callback_url", app.CallbackURL)
				d.Set("application_id", app.ApplicationID)
				d.Set("application_name", app.ApplicationName)
				d.Set("confidential", app.Confidential)
				d.SetId(strconv.Itoa(app.ID))
				return nil
			}
		}
		appsLen = len(apps)
		page = page + 1
	}

	log.Printf("[DEBUG] failed to read gitlab application %s", appID)
	d.SetId("")
	return nil
}

func resourceGitlabApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] update gitlab application %s", d.Id())

	return resourceGitlabApplicationRead(d, meta)
}

func resourceGitlabApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	log.Printf("[DEBUG] Delete gitlab application %s", d.Id())

	var err error

	appID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("failed to get clusterId: %v", err)
	}

	_, err = client.Applications.DeleteApplication(appID)
	return err
}
