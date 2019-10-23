// https://docs.gitlab.com/ee/api/merge_request_approvals.html#create-project-level-rule
package gitlab

import (
	"fmt"
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

func resourceGitLabProjectApprovalRuleSetToState(d *schema.ResourceData, rule *gitlab.ProjectApprovalRule, project string) error {
	d.SetId(buildTwoPartID(gitlab.String(project), gitlab.String(strconv.Itoa(rule.ID))))
	d.Set("project", project)
	d.Set("name", rule.Name)
	d.Set("approvals_required", rule.ApprovalsRequired)

	// if err = d.Set("user_ids", rule.Users); err != nil {
	// if err := d.Set("user_ids", flattenApprovalRuleUserIDs(rule.Users)); err != nil {
	// 	return err
	// }
	d.Set("user_ids", flattenApprovalRuleUserIDs(rule.Users))

	// if err = d.Set("group_ids", rule.Groups); err != nil {
	// if err := d.Set("user_ids", flattenApprovalRuleUserIDs(rule.Users)); err != nil {
	// 	return err
	// }
	d.Set("group_ids", flattenApprovalRuleGroupIDs(rule.Groups))

	return nil
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
		log.Printf("[DEBUG] Project %s rule created %d", project, rule.ID)
		d.SetId(buildTwoPartID(gitlab.String(project), gitlab.String(strconv.Itoa(rule.ID))))

		return resourceGitlabProjectApprovalRuleRead(d, meta)
	}

	return err
}

func resourceGitlabProjectApprovalRuleRead(d *schema.ResourceData, meta interface{}) error {
	var err error

	log.Printf("[DEBUG] read gitlab project-level rule %s", d.Id())

	rule, err := getApprovalRuleByID(meta, d.Id())
	if err != nil {
		return err
	}

	project, _, err := parseTwoPartID(d.Id())
	if err != nil {
		return err
	}

	return resourceGitLabProjectApprovalRuleSetToState(d, rule, project)
}

func resourceGitlabProjectApprovalRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	var err error

	if project, ruleID, err := parseTwoPartID(d.Id()); err == nil {
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

	if project, ruleID, err := parseTwoPartID(d.Id()); err == nil {
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

func getApprovalRuleByID(meta interface{}, id string) (*gitlab.ProjectApprovalRule, error) {
	var rules []*gitlab.ProjectApprovalRule
	var err error

	projectID, ruleIDStr, err := parseTwoPartID(id)
	fmt.Println(ruleIDStr)
	fmt.Println(projectID)
	if err != nil {
		return nil, err
	}
	ruleIDInt, err := strconv.Atoi(ruleIDStr)
	if err != nil {
		return nil, err
	}
	fmt.Println(ruleIDInt)

	client := meta.(*gitlab.Client)
	if rules, _, err = client.Projects.GetProjectApprovalRules(projectID); err == nil {
		for _, r := range rules {
			if r.ID == ruleIDInt {
				log.Printf("[DEBUG] found project-leve rule %+v", r)
				return r, nil
			}
		}
	}

	return nil, err
}

func flattenApprovalRuleUserIDs(users []*gitlab.BasicUser) []interface{} {
	// if len(users) < 1 {
	// 	return nil
	// }

	m := []interface{}{}
	for _, user := range users {
		m = append(m, user.ID)
	}

	return m
}

func flattenApprovalRuleGroupIDs(groups []*gitlab.Group) []interface{} {
	// if len(groups) < 1 {
	// 	return nil
	// }

	m := []interface{}{}
	for _, group := range groups {
		m = append(m, group.ID)
	}

	return m
}

func expandApproverIds(ids interface{}, hasItems bool) []int {
	m := []int{}
	if hasItems {
		for _, id := range ids.(*schema.Set).List() {
			m = append(m, id.(int))
		}
	}

	return m
}
