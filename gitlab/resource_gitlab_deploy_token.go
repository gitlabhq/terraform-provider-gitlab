package gitlab

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Fourcast/go-gitlab"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceGitlabDeployToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabDeployTokenCreate,
		Read:   resourceGitlabDeployTokenRead,
		Delete: resourceGitlabDeployTokenDelete,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"project", "group"},
				ForceNew:     true,
			},
			"group": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"project", "group"},
				ForceNew:     true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"expires_at": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsRFC3339Time,
				ForceNew:     true,
			},
			"scopes": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"read_registry", "read_repository"}, false),
				},
			},

			"token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceGitlabDeployTokenCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project, isProject := d.GetOk("project")
	group, isGroup := d.GetOk("group")

	var expiresAt time.Time
	var err error

	if exp, ok := d.GetOk("expires_at"); ok {
		expiresAt, err = time.Parse(time.RFC3339, exp.(string))
		if err != nil {
			return fmt.Errorf("Invalid expires_at date: %v", err)
		}
	}

	scopes := stringSetToStringSlice(d.Get("scopes").(*schema.Set))

	var deployToken *gitlab.DeployToken

	if isProject {
		options := &gitlab.CreateProjectDeployTokenOptions{
			Name:      gitlab.String(d.Get("name").(string)),
			Username:  gitlab.String(d.Get("username").(string)),
			ExpiresAt: gitlab.Time(expiresAt),
			Scopes:    *scopes,
		}

		log.Printf("[DEBUG] Create GitLab deploy token %s in project %s", *options.Name, project.(string))

		deployToken, _, err = client.DeployTokens.CreateProjectDeployToken(project, options)

	} else if isGroup {
		options := &gitlab.CreateGroupDeployTokenOptions{
			Name:      gitlab.String(d.Get("name").(string)),
			Username:  gitlab.String(d.Get("username").(string)),
			ExpiresAt: gitlab.Time(expiresAt),
			Scopes:    *scopes,
		}

		log.Printf("[DEBUG] Create GitLab deploy token %s in group %s", *options.Name, group.(string))

		deployToken, _, err = client.DeployTokens.CreateGroupDeployToken(group, options)
	}

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", deployToken.ID))

	// Token is only available on creation
	d.Set("token", deployToken.Token)
	d.Set("username", deployToken.Username)

	return nil
}

func resourceGitlabDeployTokenRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project, isProject := d.GetOk("project")
	group, isGroup := d.GetOk("group")
	deployTokenID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	var deployTokens []*gitlab.DeployToken

	if isProject {
		log.Printf("[DEBUG] Read GitLab deploy token %d in project %s", deployTokenID, project.(string))
		deployTokens, _, err = client.DeployTokens.ListProjectDeployTokens(project, nil)

	} else if isGroup {
		log.Printf("[DEBUG] Read GitLab deploy token %d in group %s", deployTokenID, group.(string))
		deployTokens, _, err = client.DeployTokens.ListGroupDeployTokens(group, nil)
	}
	if err != nil {
		return err
	}

	for _, token := range deployTokens {
		if token.ID == deployTokenID {
			d.Set("name", token.Name)
			d.Set("username", token.Username)
			d.Set("expires_at", token.ExpiresAt.Format(time.RFC3339))

			for _, scope := range token.Scopes {
				if scope == "read_repository" {
					d.Set("scopes.read_repository", true)
				}

				if scope == "read_registry" {
					d.Set("scopes.read_registry", true)
				}
			}
		}
	}

	return nil
}

func resourceGitlabDeployTokenDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project, isProject := d.GetOk("project")
	group, isGroup := d.GetOk("group")
	deployTokenID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	var response *gitlab.Response

	if isProject {
		log.Printf("[DEBUG] Delete GitLab deploy token %d in project %s", deployTokenID, project.(string))
		response, err = client.DeployTokens.DeleteProjectDeployToken(project, deployTokenID)

	} else if isGroup {
		log.Printf("[DEBUG] Delete GitLab deploy token %d in group %s", deployTokenID, group.(string))
		response, err = client.DeployTokens.DeleteGroupDeployToken(group, deployTokenID)
	}
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
