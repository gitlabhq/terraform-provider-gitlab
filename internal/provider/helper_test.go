package provider

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/onsi/gomega"
	"github.com/xanzy/go-gitlab"
)

type SkipFunc = func() (bool, error)

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

	if os.Getenv(resource.EnvTfAcc) == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' set", resource.EnvTfAcc))
	}
}

// orSkipFunc accepts many skipFunc and returns "true" if any returns true.
func orSkipFunc(input ...SkipFunc) SkipFunc {
	return func() (bool, error) {
		for _, item := range input {
			result, err := item()
			if err != nil {
				return false, err
			}
			if result {
				return result, nil
			}
		}
		return false, nil
	}
}

// testAccCheckEE is a test helper that skips the current test if the GitLab version is not GitLab Enterprise.
// This is useful when the version needs to be checked during setup, before the Terraform acceptance test starts.
func testAccCheckEE(t *testing.T) {
	t.Helper()

	version, _, err := testGitlabClient.Version.GetVersion()
	if err != nil {
		t.Fatalf("could not check GitLab version: %v", err)
	}

	if !strings.HasSuffix(version.Version, "-ee") {
		t.Skipf("Test is skipped for non-Enterprise version of GitLab (was %q)", version.String())
	}
}

// testAccCurrentUser is a test helper for getting the current user of the provided client.
func testAccCurrentUser(t *testing.T) *gitlab.User {
	t.Helper()

	user, _, err := testGitlabClient.Users.CurrentUser()
	if err != nil {
		t.Fatalf("could not get current user: %v", err)
	}

	return user
}

// testAccCreateProject is a test helper for creating a project.
func testAccCreateProject(t *testing.T) *gitlab.Project {
	return testAccCreateProjectWithNamespace(t, 0)
}

// testAccCreateProjectWithNamespace is a test helper for creating a project. This method accepts a namespace to great a project
// within a group
func testAccCreateProjectWithNamespace(t *testing.T, namespaceID int) *gitlab.Project {
	t.Helper()

	options := &gitlab.CreateProjectOptions{
		Name:        gitlab.String(acctest.RandomWithPrefix("acctest")),
		Description: gitlab.String("Terraform acceptance tests"),
		// So that acceptance tests can be run in a gitlab organization with no billing.
		Visibility: gitlab.Visibility(gitlab.PublicVisibility),
		// So that a branch is created.
		InitializeWithReadme: gitlab.Bool(true),
	}

	//Apply a namespace if one is passed in.
	if namespaceID != 0 {
		options.NamespaceID = gitlab.Int(namespaceID)
	}

	project, _, err := testGitlabClient.Projects.CreateProject(options)
	if err != nil {
		t.Fatalf("could not create test project: %v", err)
	}

	t.Cleanup(func() {
		if _, err := testGitlabClient.Projects.DeleteProject(project.ID); err != nil {
			t.Fatalf("could not cleanup test project: %v", err)
		}
	})

	return project
}

// testAccCreateUsers is a test helper for creating a specified number of users.
func testAccCreateUsers(t *testing.T, n int) []*gitlab.User {
	t.Helper()

	users := make([]*gitlab.User, n)

	for i := range users {
		var err error
		username := acctest.RandomWithPrefix("acctest-user")
		users[i], _, err = testGitlabClient.Users.CreateUser(&gitlab.CreateUserOptions{
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
			if _, err := testGitlabClient.Users.DeleteUser(userID); err != nil {
				t.Fatalf("could not cleanup test user: %v", err)
			}
		})
	}

	return users
}

// testAccCreateGroups is a test helper for creating a specified number of groups.
func testAccCreateGroups(t *testing.T, n int) []*gitlab.Group {
	t.Helper()

	groups := make([]*gitlab.Group, n)

	for i := range groups {
		var err error
		name := acctest.RandomWithPrefix("acctest-group")
		groups[i], _, err = testGitlabClient.Groups.CreateGroup(&gitlab.CreateGroupOptions{
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
			if _, err := testGitlabClient.Groups.DeleteGroup(groupID); err != nil {
				t.Fatalf("could not cleanup test group: %v", err)
			}
		})
	}

	return groups
}

// testAccCreateBranches is a test helper for creating a specified number of branches.
// It assumes the project will be destroyed at the end of the test and will not cleanup created branches.
func testAccCreateBranches(t *testing.T, project *gitlab.Project, n int) []*gitlab.Branch {
	t.Helper()

	branches := make([]*gitlab.Branch, n)

	for i := range branches {
		var err error
		branches[i], _, err = testGitlabClient.Branches.CreateBranch(project.ID, &gitlab.CreateBranchOptions{
			Branch: gitlab.String(acctest.RandomWithPrefix("acctest")),
			Ref:    gitlab.String(project.DefaultBranch),
		})
		if err != nil {
			t.Fatalf("could not create test branches: %v", err)
		}
	}

	return branches
}

// testAccCreateProtectedBranches is a test helper for creating a specified number of protected branches.
// It assumes the project will be destroyed at the end of the test and will not cleanup created branches.
func testAccCreateProtectedBranches(t *testing.T, project *gitlab.Project, n int) []*gitlab.ProtectedBranch {
	t.Helper()

	branches := testAccCreateBranches(t, project, n)
	protectedBranches := make([]*gitlab.ProtectedBranch, n)

	for i := range make([]int, n) {
		var err error
		protectedBranches[i], _, err = testGitlabClient.ProtectedBranches.ProtectRepositoryBranches(project.ID, &gitlab.ProtectRepositoryBranchesOptions{
			Name: gitlab.String(branches[i].Name),
		})
		if err != nil {
			t.Fatalf("could not protect test branches: %v", err)
		}
	}

	return protectedBranches
}

// testAccAddProjectMembers is a test helper for adding users as members of a project.
// It assumes the project will be destroyed at the end of the test and will not cleanup members.
func testAccAddProjectMembers(t *testing.T, pid interface{}, users []*gitlab.User) {
	t.Helper()

	for _, user := range users {
		_, _, err := testGitlabClient.ProjectMembers.AddProjectMember(pid, &gitlab.AddProjectMemberOptions{
			UserID:      user.ID,
			AccessLevel: gitlab.AccessLevel(gitlab.DeveloperPermissions),
		})
		if err != nil {
			t.Fatalf("could not add test project member: %v", err)
		}
	}
}

func testAccCreateProjectIssues(t *testing.T, pid interface{}, n int) []*gitlab.Issue {
	t.Helper()

	dueDate := gitlab.ISOTime(time.Now().Add(time.Hour))
	var issues []*gitlab.Issue
	for i := 0; i < n; i++ {
		issue, _, err := testGitlabClient.Issues.CreateIssue(pid, &gitlab.CreateIssueOptions{
			Title:       gitlab.String(fmt.Sprintf("Issue %d", i)),
			Description: gitlab.String(fmt.Sprintf("Description %d", i)),
			DueDate:     &dueDate,
		})
		if err != nil {
			t.Fatalf("could not create test issue: %v", err)
		}
		issues = append(issues, issue)
	}
	return issues
}

// testAccAddGroupMembers is a test helper for adding users as members of a group.
// It assumes the group will be destroyed at the end of the test and will not cleanup members.
func testAccAddGroupMembers(t *testing.T, gid interface{}, users []*gitlab.User) {
	t.Helper()

	for _, user := range users {
		_, _, err := testGitlabClient.GroupMembers.AddGroupMember(gid, &gitlab.AddGroupMemberOptions{
			UserID:      gitlab.Int(user.ID),
			AccessLevel: gitlab.AccessLevel(gitlab.DeveloperPermissions),
		})
		if err != nil {
			t.Fatalf("could not add test group member: %v", err)
		}
	}
}

func testAccAddProjectMilestone(t *testing.T, pid interface{}) *gitlab.Milestone {
	t.Helper()

	milestone, _, err := testGitlabClient.Milestones.CreateMilestone(pid, &gitlab.CreateMilestoneOptions{Title: gitlab.String("Test Milestone")})
	if err != nil {
		t.Fatalf("failed to create milestone during test for project %v: %v", pid, err)
	}
	t.Cleanup(func() {
		_, err := testGitlabClient.Milestones.DeleteMilestone(pid, milestone.ID)
		if err != nil {
			t.Fatalf("failed to delete milestone %d during test for project %v: %v", milestone.ID, pid, err)
		}
	})
	return milestone
}

func testAccCreateDeployKey(t *testing.T, projectID int, options *gitlab.AddDeployKeyOptions) *gitlab.ProjectDeployKey {
	deployKey, _, err := testGitlabClient.DeployKeys.AddDeployKey(projectID, options)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if _, err := testGitlabClient.DeployKeys.DeleteDeployKey(projectID, deployKey.ID); err != nil {
			t.Fatal(err)
		}
	})

	return deployKey
}

// testAccCreateProjectEnvironment is a test helper function for creating a project environment
func testAccCreateProjectEnvironment(t *testing.T, projectID int, options *gitlab.CreateEnvironmentOptions) *gitlab.Environment {
	t.Helper()

	projectEnvironment, _, err := testGitlabClient.Environments.CreateEnvironment(projectID, options)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if projectEnvironment.State != "stopped" {
			_, err = testGitlabClient.Environments.StopEnvironment(projectID, projectEnvironment.ID)
			if err != nil {
				t.Fatal(err)
			}
		}
		if _, err := testGitlabClient.Environments.DeleteEnvironment(projectID, projectEnvironment.ID); err != nil {
			t.Fatal(err)
		}
	})

	return projectEnvironment
}

func testAccCreateProjectVariable(t *testing.T, projectID int) *gitlab.ProjectVariable {
	variable, _, err := testGitlabClient.ProjectVariables.CreateVariable(projectID, &gitlab.CreateProjectVariableOptions{
		Key:   gitlab.String(fmt.Sprintf("test_key_%d", acctest.RandInt())),
		Value: gitlab.String("test_value"),
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if _, err := testGitlabClient.ProjectVariables.RemoveVariable(projectID, variable.Key, nil); err != nil {
			t.Fatal(err)
		}
	})

	return variable
}

func testAccCreateGroupVariable(t *testing.T, groupID int) *gitlab.GroupVariable {
	variable, _, err := testGitlabClient.GroupVariables.CreateVariable(groupID, &gitlab.CreateGroupVariableOptions{
		Key:   gitlab.String(fmt.Sprintf("test_key_%d", acctest.RandInt())),
		Value: gitlab.String("test_value"),
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if _, err := testGitlabClient.GroupVariables.RemoveVariable(groupID, variable.Key, nil); err != nil {
			t.Fatal(err)
		}
	})

	return variable
}

func testAccCreateInstanceVariable(t *testing.T) *gitlab.InstanceVariable {
	variable, _, err := testGitlabClient.InstanceVariables.CreateVariable(&gitlab.CreateInstanceVariableOptions{
		Key:   gitlab.String(fmt.Sprintf("test_key_%d", acctest.RandInt())),
		Value: gitlab.String("test_value"),
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if _, err := testGitlabClient.InstanceVariables.RemoveVariable(variable.Key, nil); err != nil {
			t.Fatal(err)
		}
	})

	return variable
}

// testAccGitlabProjectContext encapsulates a GitLab client and test project to be used during an
// acceptance test.
type testAccGitlabProjectContext struct {
	t       *testing.T
	project *gitlab.Project
}

// testAccGitlabProjectStart initializes the GitLab client and creates a test project. Remember to
// call testAccGitlabProjectContext.finish() when finished with the testAccGitlabProjectContext.
func testAccGitlabProjectStart(t *testing.T) testAccGitlabProjectContext {
	if os.Getenv(resource.EnvTfAcc) == "" {
		t.Skip(fmt.Sprintf("Acceptance tests skipped unless env '%s' set", resource.EnvTfAcc))
		return testAccGitlabProjectContext{}
	}

	project, _, err := testGitlabClient.Projects.CreateProject(&gitlab.CreateProjectOptions{
		Name:        gitlab.String(acctest.RandomWithPrefix("acctest")),
		Description: gitlab.String("Terraform acceptance tests"),
		// So that acceptance tests can be run in a gitlab organization with no billing
		Visibility: gitlab.Visibility(gitlab.PublicVisibility),
	})
	if err != nil {
		t.Fatalf("could not create test project: %v", err)
	}

	t.Cleanup(func() {
		if _, err := testGitlabClient.Projects.DeleteProject(project.ID); err != nil {
			t.Fatalf("could not delete test project: %v", err)
		}
	})

	return testAccGitlabProjectContext{
		t:       t,
		project: project,
	}
}

// testCheckResourceAttrLazy works like resource.TestCheckResourceAttr, but lazy evaluates the value parameter.
// See also: resource.TestCheckResourceAttrPtr.
func testCheckResourceAttrLazy(name string, key string, value func() string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return resource.TestCheckResourceAttr(name, key, value())(s)
	}
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
