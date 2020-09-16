package gitlab

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func resourceGitlabApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabApplicationCreate,
		Read:   resourceGitlabApplicationRead,
		Delete: resourceGitlabApplicationDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"redirect_uri": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scopes": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "",
			},
			"confidential": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"application_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceGitlabApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	var err error

	var application *gitlab.Application

	options := &gitlab.CreateApplicationOptions{
		Name:         gitlab.String(d.Get("name").(string)),
		RedirectURI:  gitlab.String(d.Get("redirect_uri").(string)),
		Scopes:       gitlab.String(d.Get("scopes").(string)),
		Confidential: gitlab.Bool(d.Get("confidential").(bool)),
	}

	log.Printf("[DEBUG] Create GitLab application %s", *options.Name)

	application, _, err = client.Applications.CreateApplication(options)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", application.ID))

	d.Set("application_id", application.ApplicationID)
	// Secret is only available on creation
	d.Set("secret", application.Secret)

	return nil
}

func resourceGitlabApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	applicationID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	var applications []*gitlab.Application

	log.Printf("[DEBUG] Read GitLab application %d", applicationID)
	applications, _, err = client.Applications.ListApplications(nil)

	if err != nil {
		return err
	}

	for _, app := range applications {
		if app.ID == applicationID {
			d.Set("name", app.ApplicationName)
			d.Set("application_id", app.ApplicationID)
			d.Set("redirect_uri", app.CallbackURL)
			d.Set("confidential", app.Confidential)
		}
	}

	return nil
}

func resourceGitlabApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	applicationID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	var response *gitlab.Response

	log.Printf("[DEBUG] Delete GitLab application %d", applicationID)
	response, err = client.Applications.DeleteApplication(applicationID)

	if err != nil {
		return err
	}

	// StatusNoContent = 204
	// Success with no body
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Invalid status code returned: %s", response.Status)
	}

	return nil
}
