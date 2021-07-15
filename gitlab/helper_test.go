package gitlab

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/xanzy/go-gitlab"
)

// testAccCompareGitLabAttribute compares an attribute in two ResourceData's for
// equivalency.
func testAccCompareGitLabAttribute(attr string, expected, received *schema.ResourceData) error {
	e := expected.Get(attr)
	r := received.Get(attr)
	switch e.(type) {
	case *schema.Set:
		if !e.(*schema.Set).Equal(r) {
			return fmt.Errorf(`attribute set %s expected "%+v" received "%+v"`, attr, e, r)
		}
	default:
		// Stringify to check because of type differences
		if fmt.Sprintf("%v", e) != fmt.Sprintf("%v", r) {
			return fmt.Errorf(`attribute %s expected "%+v" received "%+v"`, attr, e, r)
		}
	}
	return nil
}

// Returns true if the acceptance test is running Gitlab EE.
// Meant to be used as SkipFunc to skip tests that work only on Gitlab CE.
func isRunningInEE() (bool, error) {
	if conn, ok := testAccProvider.Meta().(*gitlab.Client); ok {
		version, _, err := conn.Version.GetVersion()
		if err != nil {
			return false, err
		}
		if strings.Contains(version.String(), "-ee") {
			return true, nil
		}
	} else {
		return false, errors.New("Provider not initialized, unable to get GitLab connection")
	}
	return false, nil
}

// Returns true if the acceptance test is running Gitlab CE.
// Meant to be used as SkipFunc to skip tests that work only on Gitlab EE.
func isRunningInCE() (bool, error) {
	isEE, err := isRunningInEE()
	return !isEE, err
}

// testAccGitlabProjectContext encapsulates a GitLab client and test project to be used during an
// acceptance test.
type testAccGitlabProjectContext struct {
	t       *testing.T
	client  *gitlab.Client
	project *gitlab.Project
}

// finish deletes the test project. Call this when the test is finished, usually in a defer.
func (c testAccGitlabProjectContext) finish() {
	if _, err := c.client.Projects.DeleteProject(c.project.ID); err != nil {
		c.t.Fatalf("could not delete test project: %v", err)
	}
}

// testAccGitlabProjectStart initializes the GitLab client and creates a test project. Remember to
// call testAccGitlabProjectContext.finish() when finished with the testAccGitlabProjectContext.
func testAccGitlabProjectStart(t *testing.T) testAccGitlabProjectContext {
	if os.Getenv(resource.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' set", resource.TestEnvVar))
		return testAccGitlabProjectContext{}
	}

	var options []gitlab.ClientOptionFunc
	baseURL := os.Getenv("GITLAB_BASE_URL")
	if baseURL != "" {
		options = append(options, gitlab.WithBaseURL(baseURL))
	}

	client, err := gitlab.NewClient(os.Getenv("GITLAB_TOKEN"), options...)
	if err != nil {
		t.Fatal(err)
	}

	project, _, err := client.Projects.CreateProject(&gitlab.CreateProjectOptions{
		Name:        gitlab.String(acctest.RandomWithPrefix("acctest")),
		Description: gitlab.String("Terraform acceptance tests"),
		// So that acceptance tests can be run in a gitlab organization with no billing
		Visibility: gitlab.Visibility(gitlab.PublicVisibility),
	})
	if err != nil {
		t.Fatalf("could not create test project: %v", err)
	}

	return testAccGitlabProjectContext{
		t:       t,
		client:  client,
		project: project,
	}
}
