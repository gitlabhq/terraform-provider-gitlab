package gitlab

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	gitlab "github.com/xanzy/go-gitlab"
)

var gitlabProjectApprovalRuleSkipAttributes = []string{
	"project", "user_ids", "group_ids", "eligible_approvers",
}

type checkGitlabProjectApprovalRule struct {
	SkippedAttributes []string
}

func (c *checkGitlabProjectApprovalRule) Aggregate(name string, expected gitlab.ProjectApprovalRule) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{}

	received := &gitlab.ProjectApprovalRule{}
	checks = append(checks, c.Exists(name, received))

	testResource := resourceGitlabProjectApprovalRule()
	expectedData := testResource.TestResourceData()
	receivedData := testResource.TestResourceData()

	for a, v := range testResource.Schema {
		attribute := a
		attrValue := v

		checks = append(checks, func(state *terraform.State) error {
			if testAccIsSkippedAttribute(attribute, c.SkippedAttributes) {
				switch attribute {
				case "project": // Cannot pre-compute value, but must be set.
					return resource.TestCheckResourceAttrSet(name, attribute)(state)
				case "eligible_approvers": // Checking length should be enough
					return resource.TestCheckResourceAttr(name, fmt.Sprintf("%v.#", attribute), "3")(state)
				case "user_ids": // Check if lists match
					return c.EqualLists(expectedData.Get(attribute), receivedData.Get(attribute))
				case "group_ids": // Check if lists match
					return c.EqualLists(expectedData.Get(attribute), receivedData.Get(attribute))
				default:
					return nil // Skipping
				}
			}
			if attrValue.Computed {
				attrDefault, err := attrValue.DefaultValue()
				if err != nil {
					return err
				}
				if attrDefault == nil {
					return nil
				}
			}
			expectedValue := getValueFromGitLabApprovalRule(&expected, attribute)

			return resource.TestCheckResourceAttr(name, attribute, expectedValue)(state)
		})
	}

	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func (c *checkGitlabProjectApprovalRule) EqualLists(expected, received interface{}) error {
	if expected.(*schema.Set).Len() != received.(*schema.Set).Len() {
		return fmt.Errorf("Expected set count (%d) but received (%d)", expected.(*schema.Set).Len(), received.(*schema.Set).Len())
	}

	return nil
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

// getValueFromGitLabApprovalRule is a mapping to help do a lookup from a rule
// object via string. Used for testing.
func getValueFromGitLabApprovalRule(rule *gitlab.ProjectApprovalRule, key string) string {
	var ret interface{}

	switch key {
	case "name":
		ret = rule.Name
	case "approvals_required":
		ret = rule.ApprovalsRequired
	case "group_ids":
		ret = flattenApprovalRuleGroupIDs(rule.Groups)
	case "user_ids":
		ret = flattenApprovalRuleUserIDs(rule.Users)
	case "eligible_approvers":
		ret = flattenApprovalRuleUserIDs(rule.EligibleApprovers)
	default:
		//pass
	}

	return fmt.Sprintf("%v", ret)
}

type gitlabProjectApprovalRuleFixtures struct {
	ResourceName string
	RandomInt    int
}

func (g *gitlabProjectApprovalRuleFixtures) newRule() gitlab.ProjectApprovalRule {
	fooUser := &gitlab.BasicUser{
		Name:      fmt.Sprintf("foo %d", g.RandomInt),
		Username:  fmt.Sprintf("partest-%d", g.RandomInt),
		State:     "active",
		AvatarURL: "https://www.gravatar.com/avatar/0?s=80&d=identicon",
		WebURL:    "http://localhost/foo",
	}
	barUser := &gitlab.BasicUser{
		Name:      fmt.Sprintf("bar %d", g.RandomInt),
		Username:  fmt.Sprintf("partest2-%d", g.RandomInt),
		State:     "active",
		AvatarURL: "https://www.gravatar.com/avatar/0?s=80&d=identicon",
		WebURL:    "http://localhost/foo",
	}
	barGroup := &gitlab.Group{
		Name:                 fmt.Sprintf("bar-name-%d", g.RandomInt),
		FullName:             fmt.Sprintf("bar-name-%d", g.RandomInt),
		Path:                 fmt.Sprintf("bar-path-%d", g.RandomInt),
		FullPath:             fmt.Sprintf("bar-path-%d", g.RandomInt),
		Description:          "Terraform acceptance tests - Approval Rule",
		Visibility:           gitlab.Visibility(gitlab.PublicVisibility),
		LFSEnabled:           false,
		AvatarURL:            "",
		WebURL:               fmt.Sprintf("http://localhost/groups/%d", g.RandomInt),
		RequestAccessEnabled: false,
	}

	return gitlab.ProjectApprovalRule{
		Name:              g.getName(""),
		RuleType:          "regular",
		ApprovalsRequired: 3,
		EligibleApprovers: []*gitlab.BasicUser{fooUser, barUser},
		Users:             []*gitlab.BasicUser{fooUser},
		Groups:            []*gitlab.Group{barGroup},
	}
}

func (g *gitlabProjectApprovalRuleFixtures) getName(id string) string {
	return fmt.Sprintf("foo-%s-%d", id, g.RandomInt)
}

func (g *gitlabProjectApprovalRuleFixtures) getResourceName() string {
	return fmt.Sprintf("%s.%s", g.ResourceName, "foo")
}

func (g *gitlabProjectApprovalRuleFixtures) resourceConfig() string {
	return `
resource "gitlab_project_approval_rule" "foo" {
	project            = gitlab_project.foo.id
	name               = "%s"
	approvals_required = %d
	user_ids           = [gitlab_user.foo.id]
	group_ids          = [gitlab_group.bar.id]
}

resource "gitlab_project" "foo" {
	name             = "foo-%d"
  namespace_id     = "${gitlab_group.foo.id}"
	description      = "Terraform acceptance test - Approval Rule"
	visibility_level = "public"
}

resource "gitlab_group" "foo" {
	name = "foo-name-%d"
	path = "foo-path-%d"
	description = "Terraform acceptance tests - Approval Rule"

	# So that acceptance tests can be run in a gitlab organization
	# with no billing
	visibility_level = "public"
}

resource "gitlab_user" "foo" {
	name             = "foo %d"
	username         = "partest-%d"
	password         = "partest-%dtt"
	email            = "partest-%d@ssss.com"
	is_admin         = false
	projects_limit   = 0
	can_create_group = false
	is_external      = false
}

resource "gitlab_group_membership" "foo" {
  group_id     = gitlab_group.foo.id
  user_id      = gitlab_user.foo.id
  access_level = "developer"
  expires_at   = "2024-12-31"
}

resource "gitlab_group" "bar" {
	name = "bar-name-%d"
	path = "bar-path-%d"
	description = "Terraform acceptance tests - Approval Rule"

	# So that acceptance tests can be run in a gitlab organization
	# with no billing
	visibility_level = "public"
}

resource "gitlab_user" "bar" {
	name             = "bar %d"
	username         = "partest2-%d"
	password         = "partest2-%dtt"
	email            = "partest2-%d@ssss.com"
	is_admin         = false
	projects_limit   = 0
	can_create_group = false
	is_external      = false
}

resource "gitlab_group_membership" "bar" {
  group_id     = gitlab_group.bar.id
  user_id      = gitlab_user.bar.id
  access_level = "developer"
  expires_at   = "2024-12-31"
}
`
}

func (g *gitlabProjectApprovalRuleFixtures) createConfig(name string, approvals int) string {
	return fmt.Sprintf(g.resourceConfig(),
		g.getName(name),          // approval_rule_name
		approvals,                // approvals_required
		g.RandomInt,              // project name
		g.RandomInt, g.RandomInt, // group_foo name & path
		g.RandomInt, g.RandomInt, g.RandomInt, g.RandomInt, // user_foo name, username, password, email
		g.RandomInt, g.RandomInt, // group_bar name & path
		g.RandomInt, g.RandomInt, g.RandomInt, g.RandomInt, // user_bar name, username, password, email
	)
}

func TestAccGitLabProjectApprovalRule_basic(t *testing.T) {
	testConfig := gitlabProjectApprovalRuleFixtures{
		resourceGitLabProjectApprovalRuleName,
		acctest.RandInt(),
	}
	testChecks := checkGitlabProjectApprovalRule{gitlabProjectApprovalRuleSkipAttributes}

	defaultRule := testConfig.newRule()
	updateExpected := testConfig.newRule()
	updateExpected.Name = testConfig.getName("test")
	updateExpected.ApprovalsRequired = 1

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck: func() {
			testAccPreCheck(t)
			testGitLabLicensePreCheck(t)
		},
		CheckDestroy: testChecks.Destroy,
		Steps: []resource.TestStep{
			{ // Create Rule
				Config: testConfig.createConfig("", 3),
				Check:  testChecks.Aggregate(testConfig.getResourceName(), defaultRule),
			},
			{ // Update Rule
				Config: testConfig.createConfig("test", 1),
				Check:  testChecks.Aggregate(testConfig.getResourceName(), updateExpected),
			},
			{ // Reset Rule
				Config: testConfig.createConfig("", 3),
				Check:  testChecks.Aggregate(testConfig.getResourceName(), defaultRule),
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

	defaultRule := testConfig.newRule()

	willError := defaultRule
	willError.Name = testConfig.getName("notthename")

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck: func() {
			testAccPreCheck(t)
			testGitLabLicensePreCheck(t)
		},
		CheckDestroy: testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			{ // Create rule
				Config: testConfig.createConfig("", 3),
				Check:  testChecks.Aggregate(testConfig.getResourceName(), defaultRule),
			},
			{ // Verify that name is set by passing bad values.
				Config:      testConfig.createConfig("", 3),
				Check:       testChecks.Aggregate(testConfig.getResourceName(), willError),
				ExpectError: regexp.MustCompile(`'name'\sexpected\s.+notthename.+\sgot`),
			},
			{ // Reset
				Config: testConfig.createConfig("", 3),
				Check:  testChecks.Aggregate(testConfig.getResourceName(), defaultRule),
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
		Providers: testAccProviders,
		PreCheck: func() {
			testAccPreCheck(t)
			testGitLabLicensePreCheck(t)
		},
		CheckDestroy: testAccCheckGitlabProjectDestroy,
		Steps: []resource.TestStep{
			{ // Create Rule
				Config: testConfig.createConfig("", 3),
			},
			{ // Verify Import
				ResourceName:      testConfig.getResourceName(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
