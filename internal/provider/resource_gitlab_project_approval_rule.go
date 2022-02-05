package provider

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

// https://docs.gitlab.com/ee/api/merge_request_approvals.html#create-project-level-rule
var errApprovalRuleNotFound = errors.New("approval rule not found")

var _ = registerResource("gitlab_project_approval_rule", func() *schema.Resource {
	return &schema.Resource{
		Description: "This resource allows you to create and manage multiple approval rules for your GitLab projects. For further information on approval rules, consult the [gitlab documentation](https://docs.gitlab.com/ee/api/merge_request_approvals.html#project-level-mr-approvals).\n\n" +
			"-> This feature requires GitLab Premium.",

		CreateContext: resourceGitlabProjectApprovalRuleCreate,
		ReadContext:   resourceGitlabProjectApprovalRuleRead,
		UpdateContext: resourceGitlabProjectApprovalRuleUpdate,
		DeleteContext: resourceGitlabProjectApprovalRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"project": {
				Description: "The name or id of the project to add the approval rules.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
			},
			"name": {
				Description: "The name of the approval rule.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"approvals_required": {
				Description: "The number of approvals required for this rule.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"user_ids": {
				Description: "A list of specific User IDs to add to the list of approvers.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Set:         schema.HashInt,
			},
			"group_ids": {
				Description: "A list of group IDs whose members can approve of the merge request.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Set:         schema.HashInt,
			},
			"protected_branch_ids": {
				Description: "A list of protected branch IDs (not branch names) for which the rule applies.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Set:         schema.HashInt,
			},
		},
	}
})

func resourceGitlabProjectApprovalRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	options := gitlab.CreateProjectLevelRuleOptions{
		Name:               gitlab.String(d.Get("name").(string)),
		ApprovalsRequired:  gitlab.Int(d.Get("approvals_required").(int)),
		UserIDs:            expandApproverIds(d.Get("user_ids")),
		GroupIDs:           expandApproverIds(d.Get("group_ids")),
		ProtectedBranchIDs: expandProtectedBranchIDs(d.Get("protected_branch_ids")),
	}

	project := d.Get("project").(string)

	log.Printf("[DEBUG] Project %s create gitlab project-level rule %+v", project, options)

	client := meta.(*gitlab.Client)

	rule, _, err := client.Projects.CreateProjectApprovalRule(project, &options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	ruleIDString := strconv.Itoa(rule.ID)

	d.SetId(buildTwoPartID(&project, &ruleIDString))

	return resourceGitlabProjectApprovalRuleRead(ctx, d, meta)
}

func resourceGitlabProjectApprovalRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] read gitlab project-level rule %s", d.Id())

	projectID, _, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("project", projectID)

	rule, err := getApprovalRuleByID(ctx, meta.(*gitlab.Client), d.Id())
	if err != nil {
		if errors.Is(err, errApprovalRuleNotFound) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("name", rule.Name)
	d.Set("approvals_required", rule.ApprovalsRequired)

	if err := d.Set("group_ids", flattenApprovalRuleGroupIDs(rule.Groups)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("user_ids", flattenApprovalRuleUserIDs(rule.Users)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("protected_branch_ids", flattenProtectedBranchIDs(rule.ProtectedBranches)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGitlabProjectApprovalRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectID, ruleID, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	ruleIDInt, err := strconv.Atoi(ruleID)
	if err != nil {
		return diag.FromErr(err)
	}

	options := gitlab.UpdateProjectLevelRuleOptions{
		Name:               gitlab.String(d.Get("name").(string)),
		ApprovalsRequired:  gitlab.Int(d.Get("approvals_required").(int)),
		UserIDs:            expandApproverIds(d.Get("user_ids")),
		GroupIDs:           expandApproverIds(d.Get("group_ids")),
		ProtectedBranchIDs: expandProtectedBranchIDs(d.Get("protected_branch_ids")),
	}

	log.Printf("[DEBUG] Project %s update gitlab project-level approval rule %s", projectID, *options.Name)

	client := meta.(*gitlab.Client)

	_, _, err = client.Projects.UpdateProjectApprovalRule(projectID, ruleIDInt, &options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGitlabProjectApprovalRuleRead(ctx, d, meta)
}

func resourceGitlabProjectApprovalRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	project, ruleID, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	ruleIDInt, err := strconv.Atoi(ruleID)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Project %s delete gitlab project-level approval rule %d", project, ruleIDInt)

	client := meta.(*gitlab.Client)

	_, err = client.Projects.DeleteProjectApprovalRule(project, ruleIDInt, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// getApprovalRuleByID checks the list of rules and finds the one that matches our rule ID.
func getApprovalRuleByID(ctx context.Context, client *gitlab.Client, id string) (*gitlab.ProjectApprovalRule, error) {
	projectID, ruleID, err := parseTwoPartID(id)
	if err != nil {
		return nil, err
	}

	ruleIDInt, err := strconv.Atoi(ruleID)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] read approval rules for project %s", projectID)

	rules, _, err := client.Projects.GetProjectApprovalRules(projectID, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	for _, r := range rules {
		if r.ID == ruleIDInt {
			log.Printf("[DEBUG] found project-level rule %+v", r)
			return r, nil
		}
	}

	return nil, errApprovalRuleNotFound
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

func flattenProtectedBranchIDs(protectedBranches []*gitlab.ProtectedBranch) []int {
	var protectedBranchIDs []int

	for _, protectedBranch := range protectedBranches {
		protectedBranchIDs = append(protectedBranchIDs, protectedBranch.ID)
	}

	return protectedBranchIDs
}

// expandApproverIds Expands an interface into a list of ints to read from state.
func expandApproverIds(ids interface{}) *[]int {
	var approverIDs []int

	for _, id := range ids.(*schema.Set).List() {
		approverIDs = append(approverIDs, id.(int))
	}

	return &approverIDs
}

func expandProtectedBranchIDs(ids interface{}) *[]int {
	var protectedBranchIDs []int

	for _, id := range ids.(*schema.Set).List() {
		protectedBranchIDs = append(protectedBranchIDs, id.(int))
	}

	return &protectedBranchIDs
}
