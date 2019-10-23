package gitlab

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

// https://docs.gitlab.com/ee/api/merge_request_approvals.html#create-project-level-rule

var resourceGitLabProjectApprovalRuleName = "gitlab_project_approval_rule"

var resourceGitLabProjectApprovalRuleSchema = map[string]*schema.Schema{
	"project": {
		Type:     schema.TypeString,
		ForceNew: true,
		Required: true,
	},
	"name": {
		Type:     schema.TypeString,
		ForceNew: true,
		Required: true,
	},
	"approvals_required": {
		Type:     schema.TypeInt,
		Required: true,
	},
	"user_ids": {
		Type:     schema.TypeSet,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeInt},
		Set:      schema.HashInt,
	},
	"group_ids": {
		Type:     schema.TypeSet,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeInt},
		Set:      schema.HashInt,
	},
	"eligible_approvers": {
		Type:     schema.TypeSet,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeInt},
		Set:      schema.HashInt,
	},
}

func resourceGitlabProjectApprovalRule() *schema.Resource {
	return &schema.Resource{
		Schema: resourceGitLabProjectApprovalRuleSchema,
		Create: resourceGitlabProjectApprovalRuleCreate,
		Read:   resourceGitlabProjectApprovalRuleRead,
		Update: resourceGitlabProjectApprovalRuleUpdate,
		Delete: resourceGitlabProjectApprovalRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceGitlabProjectApprovalRuleCreate(d *schema.ResourceData, meta interface{}) error {
	var err error

	options := gitlab.CreateProjectLevelRuleOptions{
		Name:              gitlab.String(d.Get("name").(string)),
		ApprovalsRequired: gitlab.Int(d.Get("approvals_required").(int)),
		UserIDs:           expandApproverIds(d.GetOk("user_ids")),
		GroupIDs:          expandApproverIds(d.GetOk("group_ids")),
	}
	project := d.Get("project").(string)

	log.Printf("[DEBUG] Project %s create gitlab project-level rule %+v", project, options)

	client := meta.(*gitlab.Client)
	if rule, _, err := client.Projects.CreateProjectApprovalRule(project, &options); err == nil {
		d.SetId(createApprovalRuleID(rule.ID, project))

		return resourceGitlabProjectApprovalRuleRead(d, meta)
	}

	return err
}

func resourceGitlabProjectApprovalRuleRead(d *schema.ResourceData, meta interface{}) error {
	var err error

	log.Printf("[DEBUG] read gitlab project-level rule %s", d.Id())

	d.Partial(true)

	projectID, _, err := splitApprovalRuleID(d.Id())
	if err != nil {
		return err
	}
	d.Set("project", projectID)
	d.SetPartial("project")

	rule, err := getApprovalRuleByID(meta, d.Id())
	if err != nil || rule == nil {
		if err != nil {
			return err
		}
		return fmt.Errorf("unable to read GitLab approvel rule %s", d.Id())
	}

	d.Set("name", rule.Name)
	d.SetPartial("name")
	d.Set("approvals_required", rule.ApprovalsRequired)
	d.SetPartial("approvals_required")

	groupIDs := flattenApprovalRuleGroupIDs(rule.Groups)
	if err := d.Set("group_ids", groupIDs); err != nil {
		return err
	}
	d.SetPartial("group_ids")

	eligibleApprovers := flattenApprovalRuleUserIDs(rule.EligibleApprovers)
	if err := d.Set("eligible_approvers", eligibleApprovers); err != nil {
		return err
	}
	d.SetPartial("eligible_approvers")

	userIDs, err := getProjectApprovalRuleUserIDs(meta.(*gitlab.Client), d.Id())
	if err != nil {
		return err
	}

	if err := d.Set("user_ids", userIDs); err != nil {
		return err
	}

	d.Partial(false)

	return nil
}

func resourceGitlabProjectApprovalRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	var err error

	if project, ruleID, err := splitApprovalRuleID(d.Id()); err == nil {
		options := gitlab.UpdateProjectLevelRuleOptions{
			Name:              gitlab.String(d.Get("name").(string)),
			ApprovalsRequired: gitlab.Int(d.Get("approvals_required").(int)),
			UserIDs:           expandApproverIds(d.GetOk("user_ids")),
			GroupIDs:          expandApproverIds(d.GetOk("group_ids")),
		}

		ruleIDInt, err := strconv.Atoi(ruleID)
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] Project %s update gitlab project-level approval rule %s", project, *options.Name)

		client := meta.(*gitlab.Client)
		if _, _, err = client.Projects.UpdateProjectApprovalRule(project, ruleIDInt, &options); err == nil {
			return resourceGitlabProjectVariableRead(d, meta)
		}
	}

	return err
}

func resourceGitlabProjectApprovalRuleDelete(d *schema.ResourceData, meta interface{}) error {
	var err error

	if project, ruleID, err := splitApprovalRuleID(d.Id()); err == nil {
		ruleIDInt, err := strconv.Atoi(ruleID)
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] Project %s delete gitlab project-level approval rule %d", project, ruleIDInt)

		client := meta.(*gitlab.Client)
		_, err = client.Projects.DeleteProjectApprovalRule(project, ruleIDInt)
	}

	return err
}

// getApprovalRuleByID checks the list of rules and finds the one that matches our rule ID.
func getApprovalRuleByID(meta interface{}, id string) (*gitlab.ProjectApprovalRule, error) {
	var rules []*gitlab.ProjectApprovalRule
	var err error

	projectID, ruleIDStr, err := splitApprovalRuleID(id)
	if err != nil {
		return nil, err
	}

	client := meta.(*gitlab.Client)
	if rules, _, err = client.Projects.GetProjectApprovalRules(projectID); err == nil {
		var ruleIDInt int

		if ruleIDInt, err = strconv.Atoi(ruleIDStr); err != nil {
			return nil, err
		}

		for _, r := range rules {
			if r.ID == ruleIDInt {
				log.Printf("[DEBUG] found project-level rule %+v", r)
				return r, nil
			}
		}
	}

	return nil, err
}

// createApprovalRuleID creates an ID of two parts: projectID:ruleID
func createApprovalRuleID(ruleID int, projectID string) string {
	return buildTwoPartID(gitlab.String(projectID), gitlab.String(strconv.Itoa(ruleID)))
}

// splitApprovalRuleID splits an approval rule into two parts, projectID:ruleID.
// An error is returned if one occurs.
func splitApprovalRuleID(ruleID string) (string, string, error) {
	return parseTwoPartID(ruleID)
}

// flattenApprovalRuleUserIDs flattens a list of approval user ids into a list
// of ints for storage in state.
func flattenApprovalRuleUserIDs(users []*gitlab.BasicUser) []interface{} {
	m := []interface{}{}
	for _, user := range users {
		m = append(m, user.ID)
	}

	return m
}

// flattenApprovalRuleGroupIDs flattens a list of approval group ids into a list
// of ints for storage in state.
func flattenApprovalRuleGroupIDs(groups []*gitlab.Group) []interface{} {
	m := []interface{}{}
	for _, group := range groups {
		m = append(m, group.ID)
	}

	return m
}

// expandApproverIds Expands an interface into a list of ints to read from state.
func expandApproverIds(ids interface{}, hasItems bool) []int {
	m := []int{}
	if hasItems {
		for _, id := range ids.(*schema.Set).List() {
			m = append(m, id.(int))
		}
	}

	return m
}

// getProjectApprovalRuleUserIDs generates a list of approvers for a rule that
// could not be found in the project or in any attached groups and returns a
// list containing only the user IDs not found as developer <= higher in the
// project, or developer <= in any of the rule groups.
//
// This method attempts to work around a bug in GitLab where there are instances
// where the `users` attribute in the response will be empty under certain
// conditions.
//
// Reference: https://gitlab.com/gitlab-org/gitlab/issues/35008
func getProjectApprovalRuleUserIDs(client *gitlab.Client, ruleID string) ([]int, error) {
	log.Print("[DEBUG] Looking up project approval rule User IDs as a workaround for bug in GitLab.")

	userIDs := newMemberList([]int{})

	rule, err := getApprovalRuleByID(client, ruleID)
	if err != nil {
		return userIDs.List(), err
	}
	users := []int{}

	if len(rule.Users) > 0 {
		for _, i := range flattenApprovalRuleUserIDs(rule.Users) {
			users = append(users, i.(int))
		}

		return users, nil
	}

	for _, i := range flattenApprovalRuleUserIDs(rule.EligibleApprovers) {
		users = append(users, i.(int))
	}

	projectID, _, err := splitApprovalRuleID(ruleID)
	if err != nil {
		return users, err
	}

	if users, err := getApproversNotInProject(client.ProjectMembers, users, projectID); err != nil {
		return users, err
	}

	groups := []int{}
	for _, i := range flattenApprovalRuleGroupIDs(rule.Groups) {
		groups = append(groups, i.(int))
	}

	if users, err := getApproversNotInGroups(client.Groups, users, groups); err != nil {
		return users, err
	}

	return users, err
}

// getApproversNotInProject returns the difference between `userIDs` and
// projectMembers to generate a list of ints containing only those not
// >= developer level.
func getApproversNotInProject(client *gitlab.ProjectMembersService, userIDs []int, projectID string) ([]int, error) {
	var err error

	memberIDs := newMemberList([]int{})
	var members []*gitlab.ProjectMember
	if members, _, err = client.ListProjectMembers(projectID, nil, nil); err != nil {
		return memberIDs.List(), err
	}

	if len(members) != 0 {
		for _, member := range members {
			if member.AccessLevel >= gitlab.DeveloperPermissions {
				memberIDs.Add(member.ID)
			}
		}
	}

	diff := newMemberList(userIDs).Difference(memberIDs.List())

	return diff, nil
}

// getApproversNotInGroups returns the difference between the users in the
// userIDs list and all of the eligible approvers found in the groups
// identified by the groupIDs list.
func getApproversNotInGroups(client *gitlab.GroupsService, userIDs []int, groupIDs []int) ([]int, error) {
	var err error

	groupMemberIds := newMemberList([]int{})

	for _, gid := range groupIDs {
		var members []*gitlab.GroupMember

		if members, _, err = client.ListGroupMembers(gid, nil, nil); err != nil {
			return nil, err
		}

		if len(members) != 0 {
			for _, gmember := range members {
				if gmember.AccessLevel >= 30 {
					groupMemberIds.Add(gmember.ID)
				}
			}
		}
	}

	diff := newMemberList(userIDs).Difference(groupMemberIds.List())

	return diff, nil
}

// memberList is used to provide collection features when working with a list of
// ints in member lists.
type memberList struct {
	ids []int
}

func newMemberList(ids []int) *memberList {
	return &memberList{
		ids: ids,
	}
}

// List returns the list of member ids.
func (m *memberList) List() []int {
	return m.ids
}

// Count is the number of members in the list.
func (m *memberList) Count() int {
	return len(m.ids)
}

// Add a new member to the id list.
func (m *memberList) Add(ids ...int) {
	for _, id := range ids {
		if !m.Contains(id) {
			m.ids = append(m.ids, id)
		}
	}
}

// Remove item from list, if it exists.
func (m *memberList) Remove(items ...int) {
	for _, item := range items {
		for index, i := range m.ids {
			if i == item {
				ids := m.List()
				ids[index] = ids[m.Count()-1]
				m.ids = ids[:m.Count()-1]
			}
		}
	}
}

// Difference returns a list of items containing only the items not found in the
// argument. ( s - t )
func (m *memberList) Difference(ids []int) []int {
	out := newMemberList(m.ids)
	out.Remove(ids...)

	return out.List()
}

// Contains returns true if the ID is in the list, false if not.
func (m *memberList) Contains(id int) bool {
	for _, i := range m.ids {
		if i == id {
			return true
		}
	}

	return false
}
