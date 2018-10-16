package gitlab

import "github.com/hashicorp/terraform/helper/schema"

// resourceGitlabServiceSlackImport is a StateFunc used to import GitLab Slack
// service by its project id.
func resourceGitlabServiceSlackImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// Actually, integration services have their own internal ids.
	// But, there is no way in GitLab API to reference any integration service
	// by its id. Instead, API for services works with project id.
	// Thus, this workaround allows us to import gitlab service
	// by the project id.

	d.Set("project", d.Id())
	return []*schema.ResourceData{d}, nil
}
