package gitlab

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

// https://docs.gitlab.com/ee/api/merge_request_approvals.html#create-project-level-rule

func resourceGitlabProjectApprovalRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabProjectApprovalRuleCreate,
		Read:   resourceGitlabProjectApprovalRuleRead,
		Update: resourceGitlabProjectApprovalRuleUpdate,
		Delete: resourceGitlabProjectApprovalRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
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
		},
	}
}

func resourceGitlabProjectApprovalRuleCreate(d *schema.ResourceData, meta interface{}) error {
	options := gitlab.CreateProjectLevelRuleOptions{
		Name:              gitlab.String(d.Get("name").(string)),
		ApprovalsRequired: gitlab.Int(d.Get("approvals_required").(int)),
		UserIDs:           expandApproverIds(d.GetOk("user_ids")),
		GroupIDs:          expandApproverIds(d.GetOk("group_ids")),
	}
	project := d.Get("project").(string)

	log.Printf("[DEBUG] Project %s create gitlab project-level rule %+v", project, options)

	client := meta.(*gitlab.Client)

	rule, _, err := client.Projects.CreateProjectApprovalRule(project, &options)
	if err != nil {
		return err
	}

	d.SetId(createApprovalRuleID(rule.ID, project))

	return resourceGitlabProjectApprovalRuleRead(d, meta)
}

func resourceGitlabProjectApprovalRuleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] read gitlab project-level rule %s", d.Id())

	d.Partial(true)

	projectID, _, err := splitApprovalRuleID(d.Id())
	if err != nil {
		return err
	}
	d.Set("project", projectID)
	d.SetPartial("project")

	rule, err := getApprovalRuleByID(meta.(*gitlab.Client), d.Id())
	if err != nil {
		return err
	}

	d.Set("name", rule.Name)
	d.SetPartial("name")

	d.Set("approvals_required", rule.ApprovalsRequired)
	d.SetPartial("approvals_required")

	if err := d.Set("group_ids", flattenApprovalRuleGroupIDs(rule.Groups)); err != nil {
		return err
	}
	d.SetPartial("group_ids")

	if err := d.Set("eligible_approvers", flattenApprovalRuleUserIDs(rule.EligibleApprovers)); err != nil {
		return err
	}
	d.SetPartial("eligible_approvers")

	if err := d.Set("user_ids", flattenApprovalRuleUserIDs(rule.Users)); err != nil {
		return err
	}
	d.SetPartial("user_ids")

	d.Partial(false)

	return nil
}

func resourceGitlabProjectApprovalRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	projectID, ruleID, err := splitApprovalRuleID(d.Id())
	if err != nil {
		return err
	}

	ruleIDInt, err := strconv.Atoi(ruleID)
	if err != nil {
		return err
	}

	options := gitlab.UpdateProjectLevelRuleOptions{
		Name:              gitlab.String(d.Get("name").(string)),
		ApprovalsRequired: gitlab.Int(d.Get("approvals_required").(int)),
		UserIDs:           expandApproverIds(d.GetOk("user_ids")),
		GroupIDs:          expandApproverIds(d.GetOk("group_ids")),
	}

	log.Printf("[DEBUG] Project %s update gitlab project-level approval rule %s", projectID, *options.Name)

	client := meta.(*gitlab.Client)

	_, _, err = client.Projects.UpdateProjectApprovalRule(projectID, ruleIDInt, &options)
	if err != nil {
		return err
	}

	return resourceGitlabProjectApprovalRuleRead(d, meta)
}

func resourceGitlabProjectApprovalRuleDelete(d *schema.ResourceData, meta interface{}) error {
	project, ruleID, err := splitApprovalRuleID(d.Id())
	if err != nil {
		return err
	}

	ruleIDInt, err := strconv.Atoi(ruleID)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Project %s delete gitlab project-level approval rule %d", project, ruleIDInt)

	client := meta.(*gitlab.Client)

	_, err = client.Projects.DeleteProjectApprovalRule(project, ruleIDInt)
	if err != nil {
		return err
	}

	return nil
}

// getApprovalRuleByID checks the list of rules and finds the one that matches our rule ID.
func getApprovalRuleByID(client *gitlab.Client, id string) (*gitlab.ProjectApprovalRule, error) {
	projectID, ruleID, err := splitApprovalRuleID(id)
	if err != nil {
		return nil, err
	}

	ruleIDInt, err := strconv.Atoi(ruleID)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] read approval rules for project %s", projectID)

	rules, _, err := client.Projects.GetProjectApprovalRules(projectID)
	if err != nil {
		return nil, err
	}

	for _, r := range rules {
		if r.ID == ruleIDInt {
			log.Printf("[DEBUG] found project-level rule %+v", r)
			return r, nil
		}
	}

	return nil, fmt.Errorf("unable to find GitLab approval rule %d", ruleIDInt)
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
func flattenApprovalRuleUserIDs(users []*gitlab.BasicUser) []int {
	var userIDs []int

	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}

	return userIDs
}

// flattenApprovalRuleGroupIDs flattens a list of approval group ids into a list
// of ints for storage in state.
func flattenApprovalRuleGroupIDs(groups []*gitlab.Group) []int {
	var groupIDs []int

	for _, group := range groups {
		groupIDs = append(groupIDs, group.ID)
	}

	return groupIDs
}

// expandApproverIds Expands an interface into a list of ints to read from state.
func expandApproverIds(ids interface{}, hasItems bool) []int {
	var approverIDs []int

	if hasItems {
		for _, id := range ids.(*schema.Set).List() {
			approverIDs = append(approverIDs, id.(int))
		}
	}

	return approverIDs
}
