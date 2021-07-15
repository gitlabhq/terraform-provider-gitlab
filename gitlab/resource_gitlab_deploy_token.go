package gitlab

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/xanzy/go-gitlab"
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
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.IsRFC3339Time,
				DiffSuppressFunc: expiresAtSuppressFunc,
				ForceNew:         true,
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

func expiresAtSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	oldDate, oldDateErr := time.Parse(time.RFC3339, old)
	newDate, newDateErr := time.Parse(time.RFC3339, new)
	if oldDateErr != nil || newDateErr != nil {
		return false
	}
	return oldDate == newDate
}

func resourceGitlabDeployTokenCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project, isProject := d.GetOk("project")
	group, isGroup := d.GetOk("group")

	var expiresAt *time.Time
	var err error

	if exp, ok := d.GetOk("expires_at"); ok {
		parsedExpiresAt, err := time.Parse(time.RFC3339, exp.(string))
		expiresAt = &parsedExpiresAt
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
			ExpiresAt: expiresAt,
			Scopes:    *scopes,
		}

		log.Printf("[DEBUG] Create GitLab deploy token %s in project %s", *options.Name, project.(string))

		deployToken, _, err = client.DeployTokens.CreateProjectDeployToken(project, options)

	} else if isGroup {
		options := &gitlab.CreateGroupDeployTokenOptions{
			Name:      gitlab.String(d.Get("name").(string)),
			Username:  gitlab.String(d.Get("username").(string)),
			ExpiresAt: expiresAt,
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
	return setResourceData(d, map[string]interface{}{
		"token":    deployToken.Token,
		"username": deployToken.Username,
	})
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
			values := map[string]interface{}{
				"name":     token.Name,
				"username": token.Username,
			}

			if token.ExpiresAt != nil {
				values["expires_at"] = token.ExpiresAt.Format(time.RFC3339)
			}

			var scopes []string
			for _, scope := range token.Scopes {
				if scope == "read_repository" || scope == "read_registry" {
					scopes = append(scopes, scope)
				}
			}
			values["scopes"] = scopes

			if err := setResourceData(d, values); err != nil {
				return err
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
