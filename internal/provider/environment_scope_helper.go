package provider

import (
	"context"
	"net/url"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/xanzy/go-gitlab"
)

// withEnvironmentScopeFilter adds the environment scope filter query parameter to the URL.
// This function is supposed to be used as `gitlab.RequestOptionFunc` parameter.
// The parameter is documented in the upstream GitLab API docs:
// https://docs.gitlab.com/ee/api/project_level_variables.html#the-filter-parameter
func withEnvironmentScopeFilter(ctx context.Context, environmentScope string) gitlab.RequestOptionFunc {
	return func(req *retryablehttp.Request) error {
		*req = *req.WithContext(ctx)
		query, err := url.ParseQuery(req.Request.URL.RawQuery)
		if err != nil {
			return err
		}
		query.Set("filter[environment_scope]", environmentScope)
		req.Request.URL.RawQuery = query.Encode()
		return nil
	}
}
