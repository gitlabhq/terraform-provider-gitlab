package provider

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mitchellh/hashstructure"
	"github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_project_milestones", func() *schema.Resource {
	validMilestoneStates := []string{"active", "closed"}

	return &schema.Resource{
		Description: `The ` + "`gitlab_project_milestones`" + ` data source allows get details of a project milestones.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/milestones.html)`,

		ReadContext: dataSourceGitlabProjectMilestonesRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Description: "The ID or URL-encoded path of the project owned by the authenticated user.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"iids": {
				Description: "Return only the milestones having the given `iid` (Note: ignored if `include_parent_milestones` is set as `true`).",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Optional:    true,
			},
			"state": {
				Description:      "Return only `active` or `closed` milestones.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(validMilestoneStates, false)),
			},
			"title": {
				Description: "Return only the milestones having the given `title`.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"search": {
				Description: "Return only milestones with a title or description matching the provided string.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"include_parent_milestones": {
				Description: "Include group milestones from parent group and its ancestors. Introduced in GitLab 13.4.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"milestones": {
				Description: "List of milestones from a project.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: datasourceSchemaFromResourceSchema(gitlabProjectMilestoneGetSchema(), nil, nil),
				},
			},
		},
	}
})

func dataSourceGitlabProjectMilestonesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	project := d.Get("project").(string)
	options := gitlab.ListMilestonesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 20,
			Page:    1,
		},
	}

	if v, ok := d.GetOk("iids"); ok {
		options.IIDs = intSetToIntSlice(v.(*schema.Set))
	}

	if v, ok := d.GetOk("title"); ok {
		options.Title = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("state"); ok {
		options.State = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("search"); ok {
		options.Search = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("include_parent_milestones"); ok {
		options.IncludeParentMilestones = gitlab.Bool(v.(bool))
	}

	optionsHash, err := hashstructure.Hash(&options, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var milestones []*gitlab.Milestone
	for options.Page != 0 {
		paginatedMilestones, resp, err := client.Milestones.ListMilestones(project, &options, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}

		milestones = append(milestones, paginatedMilestones...)
		options.Page = resp.NextPage
	}

	log.Printf("[DEBUG] get gitlab milestones from project: %s", project)
	d.SetId(fmt.Sprintf("%s:%d", project, optionsHash))
	if err = d.Set("milestones", flattenGitlabProjectMilestones(project, milestones)); err != nil {
		return diag.Errorf("Failed to set milestones to state: %v", err)
	}

	return nil
}

func flattenGitlabProjectMilestones(project string, milestones []*gitlab.Milestone) (values []map[string]interface{}) {
	for _, milestone := range milestones {
		values = append(values, gitlabProjectMilestoneToStateMap(project, milestone))
	}
	return values
}
