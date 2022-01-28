package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		provider := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"token": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("GITLAB_TOKEN", nil),
					Description: "The OAuth2 Token, Project, Group, Personal Access Token or CI Job Token used to connect to GitLab. The OAuth method is used in this provider for authentication (using Bearer authorization token). See https://docs.gitlab.com/ee/api/#authentication for details. It may be sourced from the `GITLAB_TOKEN` environment variable.",
				},
				"base_url": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("GITLAB_BASE_URL", ""),
					Description: "This is the target GitLab base API endpoint. Providing a value is a requirement when working with GitLab CE or GitLab Enterprise e.g. `https://my.gitlab.server/api/v4/`. It is optional to provide this value and it can also be sourced from the `GITLAB_BASE_URL` environment variable. The value must end with a slash.",
					ValidateFunc: func(value interface{}, key string) (ws []string, es []error) {
						v := value.(string)
						if strings.HasSuffix(v, "/api/v3") || strings.HasSuffix(v, "/api/v3/") {
							es = append(es, fmt.Errorf("terraform-provider-gitlab does not support v3 api; please upgrade to /api/v4 in %s", v))
						}
						return
					},
				},
				"cacert_file": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: "This is a file containing the ca cert to verify the gitlab instance. This is available for use when working with GitLab CE or Gitlab Enterprise with a locally-issued or self-signed certificate chain.",
				},
				"insecure": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "When set to true this disables SSL verification of the connection to the GitLab instance.",
				},
				"client_cert": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: "File path to client certificate when GitLab instance is behind company proxy. File must contain PEM encoded data.",
				},
				"client_key": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: "File path to client key when GitLab instance is behind company proxy. File must contain PEM encoded data. Required when `client_cert` is set.",
				},
				"early_auth_check": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "(Experimental) By default the provider does a dummy request to get the current user in order to verify that the provider configuration is correct and the GitLab API is reachable. Turn it off, to skip this check. This may be useful if the GitLab instance does not yet exist and is created within the same terraform module. This is an experimental feature and may change in the future. Please make sure to always keep backups of your state.",
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
				"gitlab_managed_license":            resourceGitlabManagedLicense(),
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
				"gitlab_project_access_token":       resourceGitlabProjectAccessToken(),
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
				"gitlab_repository_file":            resourceGitLabRepositoryFile(),
			},
		}

		provider.ConfigureContextFunc = configure(version, provider)
		return provider
	}

}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		config := Config{
			Token:         d.Get("token").(string),
			BaseURL:       d.Get("base_url").(string),
			CACertFile:    d.Get("cacert_file").(string),
			Insecure:      d.Get("insecure").(bool),
			ClientCert:    d.Get("client_cert").(string),
			ClientKey:     d.Get("client_key").(string),
			EarlyAuthFail: d.Get("early_auth_check").(bool),
		}

		client, err := config.Client()
		if err != nil {
			return nil, diag.FromErr(err)
		}

		userAgent := p.UserAgent("terraform-provider-gitlab", version)
		client.UserAgent = userAgent

		return client, nil
	}
}
