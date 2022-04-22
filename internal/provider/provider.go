package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	allDataSources = make(map[string]func() *schema.Resource)
	allResources   = make(map[string]func() *schema.Resource)
)

// registerDataSource may be called during package initialization to register a new data source with
// the provider.
var registerDataSource = makeRegisterResourceFunc(allDataSources, "data source")

// registerResource may be called during package initialization to register a new resource with the
// provider.
var registerResource = makeRegisterResourceFunc(allResources, "resource")

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

			DataSourcesMap: resourceFactoriesToMap(allDataSources),
			ResourcesMap:   resourceFactoriesToMap(allResources),
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

		client, err := config.Client(ctx)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		userAgent := p.UserAgent("terraform-provider-gitlab", version)
		client.UserAgent = userAgent

		return client, nil
	}
}

func makeRegisterResourceFunc(factories map[string]func() *schema.Resource, resourceType string) func(name string, fn func() *schema.Resource) interface{} {
	// lintignore: R009 // panic() during package initialization is ok
	return func(name string, fn func() *schema.Resource) interface{} {
		if strings.ToLower(name) != name {
			panic(fmt.Sprintf("cannot register %s %q: name must be lowercase", resourceType, name))
		}

		const wantPrefix = "gitlab_"
		if !strings.HasPrefix(name, wantPrefix) {
			panic(fmt.Sprintf("cannot register %s %q: name must begin with %q", resourceType, name, wantPrefix))
		}

		if _, exists := factories[name]; exists {
			panic(fmt.Sprintf("cannot register %s %q: a %s with the same name already exists", resourceType, name, resourceType))
		}

		factories[name] = fn

		return nil
	}
}

func resourceFactoriesToMap(factories map[string]func() *schema.Resource) map[string]*schema.Resource {
	resourcesMap := make(map[string]*schema.Resource)

	for name, fn := range factories {
		resourcesMap[name] = fn()
	}

	return resourcesMap
}
