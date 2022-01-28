package provider

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/hashstructure"
	"github.com/xanzy/go-gitlab"
)

func dataSourceGitlabProjectProtectedBranches() *schema.Resource {
	return &schema.Resource{
		Description: "Provides details about all protected branches in a given project.",

		ReadContext: dataSourceGitlabProjectProtectedBranchesRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Description: "The integer or path with namespace that uniquely identifies the project.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"protected_branches": {
				Description: "A list of protected branches, as defined below.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The name of the protected branch.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"id": {
							Description: "The ID of the protected branch.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"push_access_levels":  dataSourceGitlabProjectProtectedBranchSchemaAccessLevels(),
						"merge_access_levels": dataSourceGitlabProjectProtectedBranchSchemaAccessLevels(),
						"allow_force_push": {
							Description: "Whether force push is allowed.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"code_owner_approval_required": {
							Description: "Reject code pushes that change files listed in the CODEOWNERS file.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

type stateProtectedBranch struct {
	ID                        int                            `json:"id,omitempty" mapstructure:"id,omitempty"`
	Name                      string                         `json:"name,omitempty" mapstructure:"name,omitempty"`
	PushAccessLevels          []stateBranchAccessDescription `json:"push_access_levels,omitempty" mapstructure:"push_access_levels,omitempty"`
	MergeAccessLevels         []stateBranchAccessDescription `json:"merge_access_levels,omitempty" mapstructure:"merge_access_levels,omitempty"`
	AllowForcePush            bool                           `json:"allow_force_push,omitempty" mapstructure:"allow_force_push,omitempty"`
	CodeOwnerApprovalRequired bool                           `json:"code_owner_approval_required,omitempty" mapstructure:"code_owner_approval_required,omitempty"`
}

func dataSourceGitlabProjectProtectedBranchesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	log.Printf("[INFO] Reading Gitlab protected branch")

	project := d.Get("project_id")

	projectObject, _, err := client.Projects.GetProject(project, &gitlab.GetProjectOptions{}, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	allProtectedBranches := make([]stateProtectedBranch, 0)
	totalPages := -1
	opts := &gitlab.ListProtectedBranchesOptions{}
	for opts.Page = 0; opts.Page != totalPages; opts.Page++ {
		// Get protected branch by project ID/path and branch name
		pbs, resp, err := client.ProtectedBranches.ListProtectedBranches(project, opts, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
		totalPages = resp.TotalPages
		for _, pb := range pbs {
			allProtectedBranches = append(allProtectedBranches, stateProtectedBranch{
				ID:                        pb.ID,
				Name:                      pb.Name,
				PushAccessLevels:          convertBranchAccessDescriptionsToStateBranchAccessDescriptions(pb.PushAccessLevels),
				MergeAccessLevels:         convertBranchAccessDescriptionsToStateBranchAccessDescriptions(pb.MergeAccessLevels),
				AllowForcePush:            pb.AllowForcePush,
				CodeOwnerApprovalRequired: pb.CodeOwnerApprovalRequired,
			})
		}
	}

	// lintignore:R004 // TODO: Resolve this tfproviderlint issue
	if err := d.Set("protected_branches", allProtectedBranches); err != nil {
		return diag.FromErr(err)
	}

	h, err := hashstructure.Hash(*opts, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d-%d", projectObject.ID, h))

	return nil
}
