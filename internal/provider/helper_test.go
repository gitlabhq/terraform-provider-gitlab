package provider

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/onsi/gomega"
	"github.com/xanzy/go-gitlab"
)

func init() {
	// We are using the gomega package for its matchers only, but it requires us to register a handler anyway.
	gomega.RegisterFailHandler(func(_ string, _ ...int) {
		panic("gomega fail handler should not be used") // lintignore: R009
	})
}

// testAccCompareGitLabAttribute compares an attribute in two ResourceData's for
// equivalency.
func testAccCompareGitLabAttribute(attr string, expected, received *schema.ResourceData) error {
	e := expected.Get(attr)
	r := received.Get(attr)
	switch e.(type) { // nolint // TODO: Resolve this golangci-lint issue: S1034: assigning the result of this type assertion to a variable (switch e := e.(type)) could eliminate type assertions in switch cases (gosimple)
	case *schema.Set:
		if !e.(*schema.Set).Equal(r) { // nolint // TODO: Resolve this golangci-lint issue: S1034(related information): could eliminate this type assertion (gosimple)
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
	version, _, err := testGitlabClient.Version.GetVersion()
	if err != nil {
		return false, err
	}
	if strings.Contains(version.String(), "-ee") {
		return true, nil
	}

	return false, nil
}

// Returns true if the acceptance test is running Gitlab CE.
// Meant to be used as SkipFunc to skip tests that work only on Gitlab EE.
func isRunningInCE() (bool, error) {
	isEE, err := isRunningInEE()
	return !isEE, err
}

// testAccCheck is a test helper that skips the current test if it is not an acceptance test.
func testAccCheck(t *testing.T) {
	t.Helper()

	if os.Getenv(resource.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' set", resource.TestEnvVar))
	}
}

// testAccNewClient is a test helper that initializes a gitlab.Client to use in tests.
// This is preferable to using the provider metadata, which can cause unexpected behavior and breaks encapsulation.
func testAccNewClient(t *testing.T) *gitlab.Client {
	t.Helper()

	var options []gitlab.ClientOptionFunc
	baseURL := os.Getenv("GITLAB_BASE_URL")
	if baseURL != "" {
		options = append(options, gitlab.WithBaseURL(baseURL))
	}

	client, err := gitlab.NewClient(os.Getenv("GITLAB_TOKEN"), options...)
	if err != nil {
		t.Fatalf("could not initialize test client: %v", err)
	}

	return client
}

// testAccCheckEE is a test helper that skips the current test if the GitLab version is not GitLab Enterprise.
// This is useful when the version needs to be checked during setup, before the Terraform acceptance test starts.
func testAccCheckEE(t *testing.T, client *gitlab.Client) {
	t.Helper()

	version, _, err := client.Version.GetVersion()
	if err != nil {
		t.Fatalf("could not check GitLab version: %v", err)
	}

	if !strings.HasSuffix(version.Version, "-ee") {
		t.Skipf("Test is skipped for non-Enterprise version of GitLab (was %q)", version.String())
	}
}

// testAccCurrentUser is a test helper for getting the current user of the provided client.
func testAccCurrentUser(t *testing.T, client *gitlab.Client) *gitlab.User {
	t.Helper()

	user, _, err := client.Users.CurrentUser()
	if err != nil {
		t.Fatalf("could not get current user: %v", err)
	}

	return user
}

// testAccCreateGroups is a test helper for creating a project.
func testAccCreateProject(t *testing.T, client *gitlab.Client) *gitlab.Project {
	t.Helper()

	project, _, err := client.Projects.CreateProject(&gitlab.CreateProjectOptions{
		Name:        gitlab.String(acctest.RandomWithPrefix("acctest")),
		Description: gitlab.String("Terraform acceptance tests"),
		// So that acceptance tests can be run in a gitlab organization with no billing.
		Visibility: gitlab.Visibility(gitlab.PublicVisibility),
		// So that a branch is created.
		InitializeWithReadme: gitlab.Bool(true),
	})
	if err != nil {
		t.Fatalf("could not create test project: %v", err)
	}

	t.Cleanup(func() {
		if _, err := client.Projects.DeleteProject(project.ID); err != nil {
			t.Fatalf("could not cleanup test project: %v", err)
		}
	})

	return project
}

// testAccCreateUsers is a test helper for creating a specified number of users.
func testAccCreateUsers(t *testing.T, client *gitlab.Client, n int) []*gitlab.User {
	t.Helper()

	users := make([]*gitlab.User, n)

	for i := range users {
		var err error
		username := acctest.RandomWithPrefix("acctest-user")
		users[i], _, err = client.Users.CreateUser(&gitlab.CreateUserOptions{
			Name:             gitlab.String(username),
			Username:         gitlab.String(username),
			Email:            gitlab.String(username + "@example.com"),
			Password:         gitlab.String(acctest.RandString(16)),
			SkipConfirmation: gitlab.Bool(true),
		})
		if err != nil {
			t.Fatalf("could not create test user: %v", err)
		}

		userID := users[i].ID // Needed for closure.
		t.Cleanup(func() {
			if _, err := client.Users.DeleteUser(userID); err != nil {
				t.Fatalf("could not cleanup test user: %v", err)
			}
		})
	}

	return users
}

// testAccCreateGroups is a test helper for creating a specified number of groups.
func testAccCreateGroups(t *testing.T, client *gitlab.Client, n int) []*gitlab.Group {
	t.Helper()

	groups := make([]*gitlab.Group, n)

	for i := range groups {
		var err error
		name := acctest.RandomWithPrefix("acctest-group")
		groups[i], _, err = client.Groups.CreateGroup(&gitlab.CreateGroupOptions{
			Name: gitlab.String(name),
			Path: gitlab.String(name),
			// So that acceptance tests can be run in a gitlab organization with no billing.
			Visibility: gitlab.Visibility(gitlab.PublicVisibility),
		})
		if err != nil {
			t.Fatalf("could not create test group: %v", err)
		}

		groupID := groups[i].ID // Needed for closure.
		t.Cleanup(func() {
			if _, err := client.Groups.DeleteGroup(groupID); err != nil {
				t.Fatalf("could not cleanup test group: %v", err)
			}
		})
	}

	return groups
}

// testAccCreateProtectedBranches is a test helper for creating a specified number of protected branches.
// It assumes the project will be destroyed at the end of the test and will not cleanup created branches.
func testAccCreateProtectedBranches(t *testing.T, client *gitlab.Client, project *gitlab.Project, n int) []*gitlab.ProtectedBranch {
	t.Helper()

	protectedBranches := make([]*gitlab.ProtectedBranch, n)

	for i := range protectedBranches {
		branch, _, err := client.Branches.CreateBranch(project.ID, &gitlab.CreateBranchOptions{
			Branch: gitlab.String(acctest.RandomWithPrefix("acctest")),
			Ref:    gitlab.String(project.DefaultBranch),
		})
		if err != nil {
			t.Fatalf("could not create test branch: %v", err)
		}

		protectedBranches[i], _, err = client.ProtectedBranches.ProtectRepositoryBranches(project.ID, &gitlab.ProtectRepositoryBranchesOptions{
			Name: gitlab.String(branch.Name),
		})
		if err != nil {
			t.Fatalf("could not protect test branch: %v", err)
		}
	}

	return protectedBranches
}

// testAccAddProjectMembers is a test helper for adding users as members of a project.
// It assumes the project will be destroyed at the end of the test and will not cleanup members.
func testAccAddProjectMembers(t *testing.T, client *gitlab.Client, pid interface{}, users []*gitlab.User) {
	t.Helper()

	for _, user := range users {
		_, _, err := client.ProjectMembers.AddProjectMember(pid, &gitlab.AddProjectMemberOptions{
			UserID:      user.ID,
			AccessLevel: gitlab.AccessLevel(gitlab.DeveloperPermissions),
		})
		if err != nil {
			t.Fatalf("could not add test project member: %v", err)
		}
	}
}

// testAccAddGroupMembers is a test helper for adding users as members of a group.
// It assumes the group will be destroyed at the end of the test and will not cleanup members.
func testAccAddGroupMembers(t *testing.T, client *gitlab.Client, gid interface{}, users []*gitlab.User) {
	t.Helper()

	for _, user := range users {
		_, _, err := client.GroupMembers.AddGroupMember(gid, &gitlab.AddGroupMemberOptions{
			UserID:      gitlab.Int(user.ID),
			AccessLevel: gitlab.AccessLevel(gitlab.DeveloperPermissions),
		})
		if err != nil {
			t.Fatalf("could not add test group member: %v", err)
		}
	}
}

// testAccGitlabProjectContext encapsulates a GitLab client and test project to be used during an
// acceptance test.
type testAccGitlabProjectContext struct {
	t       *testing.T
	client  *gitlab.Client
	project *gitlab.Project
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

	t.Cleanup(func() {
		if _, err := client.Projects.DeleteProject(project.ID); err != nil {
			t.Fatalf("could not delete test project: %v", err)
		}
	})

	return testAccGitlabProjectContext{
		t:       t,
		client:  client,
		project: project,
	}
}
