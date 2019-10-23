package gitlab

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	gitlab "github.com/xanzy/go-gitlab"
)

var gitlabProjectApprovalRuleSkipAttributes = []string{}

type checkGitlabProjectApprovalRule struct {
	SkippedAttributes []string
}

func (c *checkGitlabProjectApprovalRule) Aggregate(name string, expected *gitlab.ProjectApprovalRule) resource.TestCheckFunc {
	var received *gitlab.ProjectApprovalRule

	checks := []resource.TestCheckFunc{
		c.Exists(name, received),
	}

	testResource := resourceGitlabProjectApprovalRule()
	expectedData := testResource.TestResourceData()
	receivedData := testResource.TestResourceData()

	for a, v := range testResource.Schema {
		attribute := a
		attrValue := v
		checks = append(checks, func(_ *terraform.State) error {
			if testAccIsSkippedAttribute(attribute, c.SkippedAttributes) {
				return nil // skipping because we said so.
			}
			if attrValue.Computed {
				if attrDefault, err := attrValue.DefaultValue(); err == nil {
					if attrDefault == nil {
						return nil // Skipping because we have no way of pre-computing computed vars
					}
				} else {
					return err
				}

			}
			resourceGitLabProjectApprovalRuleSetToState(expectedData, expected)
			resourceGitLabProjectApprovalRuleSetToState(receivedData, received)

			return testAccCompareGitLabAttribute(attribute, expectedData, receivedData)
		})
	}

	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func (c *checkGitlabProjectApprovalRule) Exists(name string, rule *gitlab.ProjectApprovalRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var err error

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not Found: %s", name)
		}

		var pid, rid string
		if pid, rid, err = parseTwoPartID(rs.Primary.ID); err != nil {
			return err
		}

		if pid == "" || rid == "" {
			return fmt.Errorf(`project ID "%s" or rule ID "%s" isn't set`, pid, rid)
		}

		client := testAccProvider.Meta().(*gitlab.Client)
		if rules, _, err := client.Projects.GetProjectApprovalRules(pid); err == nil {
			var ruleID int
			if ruleID, err = strconv.Atoi(rid); err != nil {
				return err
			}

			if len(rules) == 0 {
				return nil
			}

			for _, r := range rules {
				if r.ID == ruleID {
					*rule = *r

					return nil
				}
			}
		}

		return err
	}
}

func (c *checkGitlabProjectApprovalRule) Destroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*gitlab.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != resourceGitLabProjectApprovalRuleName {
			continue
		}

		pid, rid, err := parseTwoPartID(rs.Primary.ID)
		if err != nil {
			return err
		}

		rules, _, err := conn.Projects.GetProjectApprovalRules(pid)
		if err == nil && len(rules) == 0 {
			return nil
		}

		var ruleID int

		if ruleID, err = strconv.Atoi(rid); err != nil {
			return err
		}

		if len(rules) != 0 {
			for _, rule := range rules {
				if rule.ID == ruleID {
					return fmt.Errorf("rule still exists")
				}
			}
		}
		return nil
	}

	return nil
}

type gitlabProjectApprovalRuleFixtures struct {
	ResourceName string
	RandomInt    int
}

func (g *gitlabProjectApprovalRuleFixtures) newRule() gitlab.ProjectApprovalRule {
	user1 := &gitlab.BasicUser{
		ID:        0,
		Name:      fmt.Sprintf("Foo %d", g.RandomInt),
		Username:  "foo",
		State:     "active",
		AvatarURL: "https://www.gravatar.com/avatar/0?s=80&d=identicon",
		WebURL:    "http://localhost/foo",
	}

	user2 := user1
	user2.Name = fmt.Sprintf("Foo 2 %d", g.RandomInt)
	user2.Username = "foo2"
	user2.WebURL = "http://localhost/foo2"

	group1 := &gitlab.Group{
		ID:                   0,
		Name:                 fmt.Sprintf("foo-name-%d", g.RandomInt),
		Path:                 fmt.Sprintf("foo-path-%d", g.RandomInt),
		Description:          "Terraform acceptance tests - Approval Rule",
		Visibility:           gitlab.Visibility(gitlab.PublicVisibility),
		LFSEnabled:           false,
		AvatarURL:            "",
		WebURL:               fmt.Sprintf("http://localhost/groups/%d", g.RandomInt),
		RequestAccessEnabled: false,
		FullName:             fmt.Sprintf("foo-name-%d", g.RandomInt),
		FullPath:             fmt.Sprintf("foo-path-%d", g.RandomInt),
	}

	return gitlab.ProjectApprovalRule{
		ID:                   0,
		Name:                 g.getName(""),
		RuleType:             "regular",
		ApprovalsRequired:    3,
		EligibleApprovers:    []*gitlab.BasicUser{user1, user2},
		Users:                []*gitlab.BasicUser{user1},
		Groups:               []*gitlab.Group{group1},
		ContainsHiddenGroups: false,
	}
}

func (g *gitlabProjectApprovalRuleFixtures) getName(id string) string {
	return fmt.Sprintf("foo-%s-%d", id, g.RandomInt)
}

func (g *gitlabProjectApprovalRuleFixtures) resourceConfig() string {
	return `
resource "gitlab_project_approval_rule" "foo" {
	project = gitlab_project.foo.id
	name = "%s"
	approvals_required = %d
	user_ids = [gitlab_user.foo.id]
	group_ids = [gitlab_group.foo.id]
}

resource "gitlab_project" "foo" {
	name = "foo-%d"
	description = "Terraform acceptance test - Approval Rule"
	visibility_level = "public"
}

resource "gitlab_user" "foo" {
	name             = "foo %d"
	username         = "partest%d"
	password         = "test%dtt"
	email            = "partest%d@ssss.com"
	is_admin         = false
	projects_limit   = 0
	can_create_group = false
	is_external      = false
}

resource "gitlab_group" "foo" {
	name = "foo-name-%d"
	path = "foo-path-%d"
	description = "Terraform acceptance tests - Approval Rule"

	# So that acceptance tests can be run in a gitlab organization
	# with no billing
	visibility_level = "public"
}

resource "gitlab_group_membership" "foo2" {
  group_id     = gitlab_group.foo.id
  user_id      = gitlab_user.foo2.id
  access_level = "developer"
  expires_at   = "2024-12-31"
}

resource "gitlab_user" "foo2" {
	name             = "foo %d"
	username         = "partest2%d"
	password         = "test%dtt"
	email            = "partest%d@ssss.com"
	is_admin         = false
	projects_limit   = 0
	can_create_group = false
	is_external      = false
}
`
}

func (g *gitlabProjectApprovalRuleFixtures) createConfig(name string, approvals int) string {
	return fmt.Sprintf(g.resourceConfig(),
		g.getName(name),                                    // name, id
		approvals,                                          // approvals_required
		g.RandomInt,                                        // project name
		g.RandomInt,                                        // group name
		g.RandomInt, g.RandomInt, g.RandomInt, g.RandomInt, // user1 name, username, password, email
		g.RandomInt, g.RandomInt, g.RandomInt, g.RandomInt, // user2 name, username, password, email
	)
}

func TestAccGitLabProjectApprovalRule_basic(t *testing.T) {
	testConfig := gitlabProjectApprovalRuleFixtures{
		resourceGitLabProjectApprovalRuleName,
		acctest.RandInt(),
	}
	testChecks := checkGitlabProjectApprovalRule{gitlabProjectApprovalRuleSkipAttributes}

	resourceName := fmt.Sprintf("%s.foo", testConfig.ResourceName)

	defaultRule := testConfig.newRule()
	updateExpected := testConfig.newRule()
	updateExpected.Name = testConfig.getName("test")
	updateExpected.ApprovalsRequired = 1

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testChecks.Destroy,
		Steps: []resource.TestStep{
			// Create Rule
			resource.TestStep{
				Config: testConfig.createConfig("", 3),
				Check: resource.ComposeTestCheckFunc(
					testChecks.Aggregate(resourceName, &defaultRule),
				),
			},
			// Update Rule
			resource.TestStep{
				Config: testConfig.createConfig("test", 1),
				Check: resource.ComposeTestCheckFunc(
					testChecks.Aggregate(resourceName, &updateExpected),
				),
			},
			// Reset Rule
			resource.TestStep{
				Config: testConfig.createConfig("", 3),
				Check: resource.ComposeTestCheckFunc(
					testChecks.Aggregate(resourceName, &defaultRule),
				),
			},
		},
	})
}

func TestAccGitLabProjectApprovalRule_willError(t *testing.T) {
	testConfig := gitlabProjectApprovalRuleFixtures{
		resourceGitLabProjectApprovalRuleName,
		acctest.RandInt(),
	}
	testChecks := checkGitlabProjectApprovalRule{gitlabProjectApprovalRuleSkipAttributes}

	resourceName := fmt.Sprintf("%s.foo", testConfig.ResourceName)

	defaultRule := testConfig.newRule()

	willError := defaultRule
	willError.Name = fmt.Sprintf("foo-%s-%d", "notthename", testConfig.RandomInt)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			// Create rule
			resource.TestStep{
				Config: testConfig.createConfig("", 3),
				Check: resource.ComposeTestCheckFunc(
					testChecks.Aggregate(resourceName, &defaultRule),
				),
			},
			// Verify that name is set by passing bad values.
			resource.TestStep{
				Config:      testConfig.createConfig("notthename", 3),
				ExpectError: regexp.MustCompile(`\sname\sexpected\s.+thisisnotthename.+\sreceived`),
				Check: resource.ComposeTestCheckFunc(
					testChecks.Aggregate(resourceName, &willError),
				),
			},
			// Reset
			resource.TestStep{
				Config: testConfig.createConfig("", 3),
				Check: resource.ComposeTestCheckFunc(
					testChecks.Aggregate(resourceName, &defaultRule),
				),
			},
		},
	})
}

func TestAccGitLabProjectApprovalRule_import(t *testing.T) {
	testConfig := gitlabProjectApprovalRuleFixtures{
		resourceGitLabProjectApprovalRuleName,
		acctest.RandInt(),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testConfig.createConfig("", 3),
			},
			resource.TestStep{
				ResourceName:      fmt.Sprintf("%s.foo", testConfig.ResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
