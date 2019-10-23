// https://docs.gitlab.com/ee/api/merge_request_approvals.html#create-project-level-rule
package gitlab

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

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
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
	},
	"group_ids": {
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
	},
}

func resourceGitlabProjectApprovalRule() *schema.Resource {
	return &schema.Resource{
		Schema: resourceGitLabProjectApprovalRuleSchema,
		Exists: resourceGitlabProjectApprovalRuleExists,
		Create: resourceGitlabProjectApprovalRuleCreate,
		Read:   resourceGitlabProjectApprovalRuleRead,
		Update: resourceGitlabProjectApprovalRuleUpdate,
		Delete: resourceGitlabProjectApprovalRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceGitLabProjectApprovalRuleSetToState(d *schema.ResourceData, rule *gitlab.ProjectApprovalRule) {
	d.Set("name", rule.Name)
	d.Set("approvals_required", rule.ApprovalsRequired)

	userIDs := []int{}
	for _, uid := range rule.Users {
		userIDs = append(userIDs, uid.ID)
	}
	d.Set("user_ids", userIDs)

	groupIDs := []int{}
	for _, gid := range rule.Groups {
		groupIDs = append(groupIDs, gid.ID)
	}
	d.Set("group_ids", groupIDs)
}

func resourceGitlabProjectApprovalRuleExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	var err error

	project, ruleID, err := parseTwoPartID(d.Id())
	if err != nil {
		return false, err
	}

	if rule, err := getApprovalRuleByID(meta, project, ruleID); err == nil {
		if rule != nil {
			return true, err
		}
	}

	return false, err
}

func resourceGitlabProjectApprovalRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	project := d.Get("project").(string)
	name := d.Get("name").(string)
	approvalsRequired := d.Get("approvals_required").(int)
	userIDs := d.Get("user_ids").([]int)
	groupIDs := d.Get("group_ids").([]int)

	log.Printf("[DEBUG] create gitlab project-level rule %s", name)
	options := gitlab.CreateProjectLevelRuleOptions{
		Name:              &name,
		ApprovalsRequired: &approvalsRequired,
		UserIDs:           userIDs,
		GroupIDs:          groupIDs,
	}

	var err error
	if rule, _, err := client.Projects.CreateProjectApprovalRule(project, &options); err == nil {
		ruleID := strconv.Itoa(rule.ID)
		d.SetId(buildTwoPartID(&project, &ruleID))

		return resourceGitlabProjectApprovalRuleRead(d, meta)
	}

	return err
}

func resourceGitlabProjectApprovalRuleRead(d *schema.ResourceData, meta interface{}) error {
	var err error

	if project, ruleID, err := parseTwoPartID(d.Id()); err == nil {
		log.Printf("[DEBUG] read gitlab project-level rule %s/%s", project, ruleID)

		if rule, err := getApprovalRuleByID(meta, project, ruleID); err == nil {
			resourceGitLabProjectApprovalRuleSetToState(d, rule)
		}
	}

	return err
}

func resourceGitlabProjectApprovalRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	var err error

	if project, ruleID, err := parseTwoPartID(d.Id()); err == nil {
		log.Printf("[DEBUG] update gitlab project-level approval rule %s/%s", project, ruleID)

		name := d.Get("name").(string)
		approvalsRequired := d.Get("approvals_required").(int)
		userIDs := d.Get("user_ids").([]int)
		groupIDs := d.Get("group_ids").([]int)

		options := gitlab.UpdateProjectLevelRuleOptions{
			Name:              &name,
			ApprovalsRequired: &approvalsRequired,
			UserIDs:           userIDs,
			GroupIDs:          groupIDs,
		}

		client := meta.(*gitlab.Client)

		ruleIDInt, err := strconv.Atoi(ruleID)
		if err != nil {
			return err
		}

		if _, _, err = client.Projects.UpdateProjectApprovalRule(project, ruleIDInt, &options); err == nil {
			return resourceGitlabProjectVariableRead(d, meta)
		}
	}

	return err
}

func resourceGitlabProjectApprovalRuleDelete(d *schema.ResourceData, meta interface{}) error {
	var err error

	if project, ruleID, err := parseTwoPartID(d.Id()); err == nil {
		log.Printf("[DEBUG] Delete gitlab project-level approval rule %s/%s", project, ruleID)

		client := meta.(*gitlab.Client)

		ruleIDInt, err := strconv.Atoi(ruleID)
		if err != nil {
			return err
		}

		_, err = client.Projects.DeleteProjectApprovalRule(project, ruleIDInt)
	}

	return err
}

func getApprovalRuleByID(meta interface{}, pid string, ruleID string) (*gitlab.ProjectApprovalRule, error) {
	var rules []*gitlab.ProjectApprovalRule
	var err error

	client := meta.(*gitlab.Client)

	if rules, _, err = client.Projects.GetProjectApprovalRules(pid); err != nil {
		return nil, err
	}

	ruleIDInt, err := strconv.Atoi(ruleID)
	if err != nil {
		return nil, err
	}

	for _, r := range rules {
		if r.ID == ruleIDInt {
			return r, err
		}
	}

	return nil, err
}
