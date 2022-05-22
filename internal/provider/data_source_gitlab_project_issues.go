package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mitchellh/hashstructure"
	"github.com/xanzy/go-gitlab"
)

var validIssueOrderByValues = []string{
	"created_at",
	"updated_at",
	"priority",
	"due_date",
	"relative_position",
	"label_priority",
	"milestone_due",
	"popularity",
	"weight",
}

var _ = registerDataSource("gitlab_project_issues", func() *schema.Resource {
	validIssueScopeValues := []string{"created_by_me", "assigned_to_me", "all"}
	validIssueSortValues := []string{"asc", "desc"}

	return &schema.Resource{
		Description: `The ` + "`gitlab_project_issues`" + ` data source allows to retrieve details about issues in a project.

**Upstream API**: [GitLab API docs](https://docs.gitlab.com/ee/api/issues.html)`,

		ReadContext: dataSourceGitlabProjectIssuesRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Description: "The name or id of the project.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"assignee_id": {
				Description: "Return issues assigned to the given user id. Mutually exclusive with assignee_username. None returns unassigned issues. Any returns issues with an assignee.",
				Type:        schema.TypeInt,
				Optional:    true,
				ConflictsWith: []string{
					"assignee_username",
				},
			},
			"not_assignee_id": {
				Description: "Return issues that do not match the assignee id.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Optional:    true,
			},
			"assignee_username": {
				Description: "Return issues assigned to the given username. Similar to assignee_id and mutually exclusive with assignee_id. In GitLab CE, the assignee_username array should only contain a single value. Otherwise, an invalid parameter error is returned.",
				Type:        schema.TypeString,
				Optional:    true,
				ConflictsWith: []string{
					"assignee_id",
				},
			},
			"author_id": {
				Description: "Return issues created by the given user id. Combine with scope=all or scope=assigned_to_me.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"not_author_id": {
				Description: "Return issues that do not match the author id.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Optional:    true,
			},
			// NOTE: not yet supported in go-gitlab.
			// "author_username": {},
			"confidential": {
				Description: "Filter confidential or public issues.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"created_after": {
				Description: "Return issues created on or after the given time. Expected in ISO 8601 format (2019-03-15T08:00:00Z)",
				Type:        schema.TypeString,
				Optional:    true,
				// NOTE: since RFC3339 is pretty much a subset of ISO8601 and actually expected by GitLab,
				//       we use it here to avoid having to parse the string ourselves.
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsRFC3339Time),
			},
			"created_before": {
				Description: "Return issues created on or before the given time. Expected in ISO 8601 format (2019-03-15T08:00:00Z)",
				Type:        schema.TypeString,
				Optional:    true,
				// NOTE: since RFC3339 is pretty much a subset of ISO8601 and actually expected by GitLab,
				//       we use it here to avoid having to parse the string ourselves.
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsRFC3339Time),
			},
			"due_date": {
				Description: "Return issues that have no due date, are overdue, or whose due date is this week, this month, or between two weeks ago and next month. Accepts: 0 (no due date), any, today, tomorrow, overdue, week, month, next_month_and_previous_two_weeks.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			// NOTE: not yet supported in go-gitlab.
			// "epic_id": {}
			"iids": {
				Description: "Return only the issues having the given iid",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Optional:    true,
			},
			"issue_type": {
				Description:      fmt.Sprintf("Filter to a given type of issue. Valid values are %s. (Introduced in GitLab 13.12)", validIssueTypes),
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(validIssueTypes, false)),
			},
			// NOTE: not yet supported in go-gitlab.
			// iteration_id: {},
			// iteration_title: {},
			"labels": {
				Description: "Return issues with labels. Issues must have all labels to be returned. None lists all issues with no labels. Any lists all issues with at least one label. No+Label (Deprecated) lists all issues with no labels. Predefined names are case-insensitive.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			"not_labels": {
				Description: "Return issues that do not match the labels.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			"milestone": {
				Description: "The milestone title. None lists all issues with no milestone. Any lists all issues that have an assigned milestone.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"not_milestone": {
				Description: "Return issues that do not match the milestone.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"my_reaction_emoji": {
				Description: "Return issues reacted by the authenticated user by the given emoji. None returns issues not given a reaction. Any returns issues given at least one reaction.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"not_my_reaction_emoji": {
				Description: "Return issues not reacted by the authenticated user by the given emoji.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			"order_by": {
				Description:      fmt.Sprintf("Return issues ordered by. Valid values are %s. Default is created_at", renderValueListForDocs(validIssueOrderByValues)),
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(validIssueOrderByValues, false)),
			},
			"scope": {
				Description:      fmt.Sprintf("Return issues for the given scope. Valid values are %s. Defaults to all.", renderValueListForDocs(validIssueScopeValues)),
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(validIssueScopeValues, false)),
			},
			"search": {
				Description: "Search project issues against their title and description",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"sort": {
				Description:      "Return issues sorted in asc or desc order. Default is desc",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(validIssueSortValues, false)),
			},
			"updated_after": {
				Description: "Return issues updated on or after the given time. Expected in ISO 8601 format (2019-03-15T08:00:00Z)",
				Type:        schema.TypeString,
				Optional:    true,
				// NOTE: since RFC3339 is pretty much a subset of ISO8601 and actually expected by GitLab,
				//       we use it here to avoid having to parse the string ourselves.
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsRFC3339Time),
			},
			"updated_before": {
				Description: "Return issues updated on or before the given time. Expected in ISO 8601 format (2019-03-15T08:00:00Z)",
				Type:        schema.TypeString,
				Optional:    true,
				// NOTE: since RFC3339 is pretty much a subset of ISO8601 and actually expected by GitLab,
				//       we use it here to avoid having to parse the string ourselves.
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsRFC3339Time),
			},
			"weight": {
				Description: "Return issues with the specified weight. None returns issues with no weight assigned. Any returns issues with a weight assigned.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"with_labels_details": {
				Description: "If true, the response returns more details for each label in labels field: :name, :color, :description, :description_html, :text_color. Default is false. description_html was introduced in GitLab 12.7",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"issues": {
				Description: "The list of issues returned by the search.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: datasourceSchemaFromResourceSchema(gitlabProjectIssueGetSchema(), nil, nil),
				},
			},
		},
	}
})

func dataSourceGitlabProjectIssuesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	project := d.Get("project").(string)
	options := gitlab.ListProjectIssuesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 20,
			Page:    1,
		},
	}

	if v, ok := d.GetOk("iids"); ok {
		options.IIDs = intSetToIntSlice(v.(*schema.Set))
	}

	if v, ok := d.GetOk("state"); ok {
		options.State = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("labels"); ok {
		gitlabLabels := gitlab.Labels(*stringSetToStringSlice(v.(*schema.Set)))
		options.Labels = &gitlabLabels
	}

	if v, ok := d.GetOk("not_labels"); ok {
		gitlabLabels := gitlab.Labels(*stringSetToStringSlice(v.(*schema.Set)))
		options.Labels = &gitlabLabels
	}

	if v, ok := d.GetOk("with_labels_details"); ok {
		options.WithLabelDetails = gitlab.Bool(v.(bool))
	}

	if v, ok := d.GetOk("milestone"); ok {
		options.Milestone = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("not_milestone"); ok {
		options.NotMilestone = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("scope"); ok {
		options.Scope = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("author_id"); ok {
		options.AuthorID = gitlab.Int(v.(int))
	}

	if v, ok := d.GetOk("not_author_id"); ok {
		options.NotAuthorID = intSetToIntSlice(v.(*schema.Set))
	}

	if v, ok := d.GetOk("assignee_id"); ok {
		options.AssigneeID = gitlab.AssigneeID(v.(int))
	}

	if v, ok := d.GetOk("not_assignee_id"); ok {
		options.NotAssigneeID = intSetToIntSlice(v.(*schema.Set))
	}

	if v, ok := d.GetOk("assignee_username"); ok {
		options.AssigneeUsername = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("my_reaction_emoji"); ok {
		options.MyReactionEmoji = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("not_my_reaction_emoji"); ok {
		options.NotMyReactionEmoji = stringSetToStringSlice(v.(*schema.Set))
	}

	if v, ok := d.GetOk("order_by"); ok {
		options.OrderBy = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("sort"); ok {
		options.Sort = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("search"); ok {
		options.Search = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("in"); ok {
		options.In = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("created_after"); ok {
		parsedCreatedAfter, err := time.Parse(time.RFC3339, v.(string))
		if err != nil {
			return diag.Errorf("failed to parse created_after: %s. It must be in valid RFC3339 format.", err)
		}
		options.CreatedAfter = gitlab.Time(parsedCreatedAfter)
	}

	if v, ok := d.GetOk("created_before"); ok {
		parsedCreatedBefore, err := time.Parse(time.RFC3339, v.(string))
		if err != nil {
			return diag.Errorf("failed to parse created_before: %s. It must be in valid RFC3339 format.", err)
		}
		options.CreatedBefore = gitlab.Time(parsedCreatedBefore)
	}

	if v, ok := d.GetOk("due_date"); ok {
		options.DueDate = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("updated_after"); ok {
		parsedUpdatedAfter, err := time.Parse(time.RFC3339, v.(string))
		if err != nil {
			return diag.Errorf("failed to parse updated_after: %s. It must be in valid RFC3339 format.", err)
		}
		options.UpdatedAfter = gitlab.Time(parsedUpdatedAfter)
	}

	if v, ok := d.GetOk("updated_before"); ok {
		parsedUpdatedBefore, err := time.Parse(time.RFC3339, v.(string))
		if err != nil {
			return diag.Errorf("failed to parse updated_before: %s. It must be in valid RFC3339 format.", err)
		}
		options.UpdatedBefore = gitlab.Time(parsedUpdatedBefore)
	}

	if v, ok := d.GetOk("confidential"); ok {
		options.Confidential = gitlab.Bool(v.(bool))
	}

	if v, ok := d.GetOk("issue_type"); ok {
		options.IssueType = gitlab.String(v.(string))
	}

	var issues []*gitlab.Issue
	for options.Page != 0 {
		paginatedIssues, resp, err := client.Issues.ListProjectIssues(project, &options, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}

		issues = append(issues, paginatedIssues...)
		options.Page = resp.NextPage
	}

	optionsHash, err := hashstructure.Hash(&options, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s-%d", project, optionsHash))
	if err = d.Set("issues", flattenGitlabProjectIssues(issues)); err != nil {
		return diag.Errorf("failed to set issues to state: %v", err)
	}

	return nil
}

func flattenGitlabProjectIssues(issues []*gitlab.Issue) (values []map[string]interface{}) {
	for _, issue := range issues {
		values = append(values, gitlabProjectIssueToStateMap(fmt.Sprintf("%d", issue.ProjectID), issue))
	}
	return values
}
