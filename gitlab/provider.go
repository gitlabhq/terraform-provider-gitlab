package gitlab

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {

	// The actual provider
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GITLAB_TOKEN", nil),
				Description: descriptions["token"],
			},
			"base_url": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("GITLAB_BASE_URL", ""),
				Description:  descriptions["base_url"],
				ValidateFunc: validateApiURLVersion,
			},
			"cacert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["cacert_file"],
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["insecure"],
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"gitlab_group":   dataSourceGitlabGroup(),
			"gitlab_project": dataSourceGitlabProject(),
			"gitlab_user":    dataSourceGitlabUser(),
			"gitlab_users":   dataSourceGitlabUsers(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"gitlab_branch_protection":  resourceGitlabBranchProtection(),
			"gitlab_tag_protection":     resourceGitlabTagProtection(),
			"gitlab_group":              resourceGitlabGroup(),
			"gitlab_project":            resourceGitlabProject(),
			"gitlab_label":              resourceGitlabLabel(),
			"gitlab_pipeline_schedule":  resourceGitlabPipelineSchedule(),
			"gitlab_pipeline_trigger":   resourceGitlabPipelineTrigger(),
			"gitlab_project_hook":       resourceGitlabProjectHook(),
			"gitlab_deploy_key":         resourceGitlabDeployKey(),
			"gitlab_user":               resourceGitlabUser(),
			"gitlab_project_membership": resourceGitlabProjectMembership(),
			"gitlab_group_membership":   resourceGitlabGroupMembership(),
			"gitlab_project_variable":   resourceGitlabProjectVariable(),
			"gitlab_group_variable":     resourceGitlabGroupVariable(),
			"gitlab_project_cluster":    resourceGitlabProjectCluster(),
			"gitlab_service_slack":      resourceGitlabServiceSlack(),
			"gitlab_service_jira":       resourceGitlabServiceJira(),
		},

		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"token": "The OAuth token used to connect to GitLab.",

		"base_url": "The GitLab Base API URL",

		"cacert_file": "A file containing the ca certificate to use in case ssl certificate is not from a standard chain",

		"insecure": "Disable SSL verification of API calls",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Token:      d.Get("token").(string),
		BaseURL:    d.Get("base_url").(string),
		CACertFile: d.Get("cacert_file").(string),
		Insecure:   d.Get("insecure").(bool),
	}

	return config.Client()
}

func validateApiURLVersion(value interface{}, key string) (ws []string, es []error) {
	v := value.(string)
	if strings.HasSuffix(v, "/api/v3") || strings.HasSuffix(v, "/api/v3/") {
		es = append(es, fmt.Errorf("terraform-gitlab-provider does not support v3 api; please upgrade to /api/v4 in %s", v))
	}
	return
}
