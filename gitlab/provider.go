package gitlab

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	provider := &schema.Provider{
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
			"client_cert": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["client_cert"],
			},
			"client_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: descriptions["client_key"],
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"gitlab_group":                      dataSourceGitlabGroup(),
			"gitlab_group_membership":           dataSourceGitlabGroupMembership(),
			"gitlab_project":                    dataSourceGitlabProject(),
			"gitlab_project_protected_branch":   dataSourceGitlabProjectProtectedBranch(),
			"gitlab_project_protected_branches": dataSourceGitlabProjectProtectedBranches(),
			"gitlab_projects":                   dataSourceGitlabProjects(),
			"gitlab_user":                       dataSourceGitlabUser(),
			"gitlab_users":                      dataSourceGitlabUsers(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"gitlab_branch_protection":          resourceGitlabBranchProtection(),
			"gitlab_tag_protection":             resourceGitlabTagProtection(),
			"gitlab_group":                      resourceGitlabGroup(),
			"gitlab_group_custom_attribute":     resourceGitlabGroupCustomAttribute(),
			"gitlab_project":                    resourceGitlabProject(),
			"gitlab_project_custom_attribute":   resourceGitlabProjectCustomAttribute(),
			"gitlab_label":                      resourceGitlabLabel(),
			"gitlab_group_label":                resourceGitlabGroupLabel(),
			"gitlab_pipeline_schedule":          resourceGitlabPipelineSchedule(),
			"gitlab_pipeline_schedule_variable": resourceGitlabPipelineScheduleVariable(),
			"gitlab_pipeline_trigger":           resourceGitlabPipelineTrigger(),
			"gitlab_project_hook":               resourceGitlabProjectHook(),
			"gitlab_deploy_key":                 resourceGitlabDeployKey(),
			"gitlab_deploy_key_enable":          resourceGitlabDeployEnableKey(),
			"gitlab_deploy_token":               resourceGitlabDeployToken(),
			"gitlab_user":                       resourceGitlabUser(),
			"gitlab_user_custom_attribute":      resourceGitlabUserCustomAttribute(),
			"gitlab_project_membership":         resourceGitlabProjectMembership(),
			"gitlab_group_membership":           resourceGitlabGroupMembership(),
			"gitlab_project_variable":           resourceGitlabProjectVariable(),
			"gitlab_group_variable":             resourceGitlabGroupVariable(),
			"gitlab_project_cluster":            resourceGitlabProjectCluster(),
			"gitlab_service_slack":              resourceGitlabServiceSlack(),
			"gitlab_service_jira":               resourceGitlabServiceJira(),
			"gitlab_service_microsoft_teams":    resourceGitlabServiceMicrosoftTeams(),
			"gitlab_service_github":             resourceGitlabServiceGithub(),
			"gitlab_service_pipelines_email":    resourceGitlabServicePipelinesEmail(),
			"gitlab_project_share_group":        resourceGitlabProjectShareGroup(),
			"gitlab_group_cluster":              resourceGitlabGroupCluster(),
			"gitlab_group_ldap_link":            resourceGitlabGroupLdapLink(),
			"gitlab_instance_cluster":           resourceGitlabInstanceCluster(),
			"gitlab_project_mirror":             resourceGitlabProjectMirror(),
			"gitlab_project_level_mr_approvals": resourceGitlabProjectLevelMRApprovals(),
			"gitlab_project_approval_rule":      resourceGitlabProjectApprovalRule(),
			"gitlab_instance_variable":          resourceGitlabInstanceVariable(),
			"gitlab_project_freeze_period":      resourceGitlabProjectFreezePeriod(),
			"gitlab_group_share_group":          resourceGitlabGroupShareGroup(),
			"gitlab_project_badge":              resourceGitlabProjectBadge(),
			"gitlab_group_badge":                resourceGitlabGroupBadge(),
		},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return providerConfigure(ctx, provider, d)
	}

	return provider
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"token": "The OAuth2 token or project/personal access token used to connect to GitLab.",

		"base_url": "The GitLab Base API URL",

		"cacert_file": "A file containing the ca certificate to use in case ssl certificate is not from a standard chain",

		"insecure": "Disable SSL verification of API calls",

		"client_cert": "File path to client certificate when GitLab instance is behind company proxy. File  must contain PEM encoded data.",

		"client_key": "File path to client key when GitLab instance is behind company proxy. File must contain PEM encoded data.",
	}
}

func providerConfigure(ctx context.Context, p *schema.Provider, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := Config{
		Token:      d.Get("token").(string),
		BaseURL:    d.Get("base_url").(string),
		CACertFile: d.Get("cacert_file").(string),
		Insecure:   d.Get("insecure").(bool),
		ClientCert: d.Get("client_cert").(string),
		ClientKey:  d.Get("client_key").(string),
	}

	client, err := config.Client()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	userAgent := p.UserAgent("terraform-provider-gitlab", "")
	client.UserAgent = userAgent

	return client, nil
}

func validateApiURLVersion(value interface{}, key string) (ws []string, es []error) {
	v := value.(string)
	if strings.HasSuffix(v, "/api/v3") || strings.HasSuffix(v, "/api/v3/") {
		es = append(es, fmt.Errorf("terraform-provider-gitlab does not support v3 api; please upgrade to /api/v4 in %s", v))
	}
	return
}
