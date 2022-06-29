package provider

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_group", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_group`" + ` resource allows to manage the lifecycle of a group.

-> On GitLab SaaS, you must use the GitLab UI to create groups without a parent group. You cannot use this provider nor the API to do this.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/groups.html)`,

		CreateContext: resourceGitlabGroupCreate,
		ReadContext:   resourceGitlabGroupRead,
		UpdateContext: resourceGitlabGroupUpdate,
		DeleteContext: resourceGitlabGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of this group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"path": {
				Description: "The path of the group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"full_path": {
				Description: "The full path of the group.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"full_name": {
				Description: "The full name of the group.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"web_url": {
				Description: "Web URL of the group.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "The description of the group.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"lfs_enabled": {
				Description: "Defaults to true. Enable/disable Large File Storage (LFS) for the projects in this group.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"default_branch_protection": {
				Description:  "Defaults to 2. See https://docs.gitlab.com/ee/api/groups.html#options-for-default_branch_protection",
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 2, 3}),
			},
			"request_access_enabled": {
				Description: "Defaults to false. Allow users to request member access.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"visibility_level": {
				Description:  "The group's visibility. Can be `private`, `internal`, or `public`.",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "internal", "public"}, true),
			},
			"share_with_group_lock": {
				Description: "Defaults to false. Prevent sharing a project with another group within this group.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"project_creation_level": {
				Description:  "Defaults to maintainer. Determine if developers can create projects in the group.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "maintainer",
				ValidateFunc: validation.StringInSlice([]string{"noone", "maintainer", "developer"}, true),
			},
			"auto_devops_enabled": {
				Description: "Defaults to false. Default to Auto DevOps pipeline for all projects within this group.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"emails_disabled": {
				Description: "Defaults to false. Disable email notifications.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"mentions_disabled": {
				Description: "Defaults to false. Disable the capability of a group from getting mentioned.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"subgroup_creation_level": {
				Description:  "Defaults to owner. Allowed to create subgroups.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "owner",
				ValidateFunc: validation.StringInSlice([]string{"owner", "maintainer"}, true),
			},
			"require_two_factor_authentication": {
				Description: "Defaults to false. Require all users in this group to setup Two-factor authentication.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"two_factor_grace_period": {
				Description: "Defaults to 48. Time before Two-factor authentication is enforced (in hours).",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     48,
			},
			"parent_id": {
				Description: "Id of the parent group (creates a nested group).",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
			},
			"runners_token": {
				Description: "The group level registration token to use during runner setup.",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
			"prevent_forking_outside_group": {
				Description: "Defaults to false. When enabled, users can not fork projects from this group to external namespaces.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
		},
	}
})

func resourceGitlabGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	options := &gitlab.CreateGroupOptions{
		Name:                 gitlab.String(d.Get("name").(string)),
		LFSEnabled:           gitlab.Bool(d.Get("lfs_enabled").(bool)),
		RequestAccessEnabled: gitlab.Bool(d.Get("request_access_enabled").(bool)),
	}

	if v, ok := d.GetOk("path"); ok {
		options.Path = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		options.Description = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("visibility_level"); ok {
		options.Visibility = stringToVisibilityLevel(v.(string))
	}

	if v, ok := d.GetOk("share_with_group_lock"); ok {
		options.ShareWithGroupLock = gitlab.Bool(v.(bool))
	}

	// nolint:staticcheck // SA1019 ignore deprecated GetOkExists
	// lintignore: XR001 // TODO: replace with alternative for GetOkExists
	if v, ok := d.GetOkExists("require_two_factor_authentication"); ok {
		options.RequireTwoFactorAuth = gitlab.Bool(v.(bool))
	}

	if v, ok := d.GetOk("two_factor_grace_period"); ok {
		options.TwoFactorGracePeriod = gitlab.Int(v.(int))
	}

	if v, ok := d.GetOk("project_creation_level"); ok {
		options.ProjectCreationLevel = stringToProjectCreationLevel(v.(string))
	}

	// nolint:staticcheck // SA1019 ignore deprecated GetOkExists
	// lintignore: XR001 // TODO: replace with alternative for GetOkExists
	if v, ok := d.GetOkExists("auto_devops_enabled"); ok {
		options.AutoDevopsEnabled = gitlab.Bool(v.(bool))
	}

	if v, ok := d.GetOk("subgroup_creation_level"); ok {
		options.SubGroupCreationLevel = stringToSubGroupCreationLevel(v.(string))
	}

	// nolint:staticcheck // SA1019 ignore deprecated GetOkExists
	// lintignore: XR001 // TODO: replace with alternative for GetOkExists
	if v, ok := d.GetOkExists("emails_disabled"); ok {
		options.EmailsDisabled = gitlab.Bool(v.(bool))
	}

	// nolint:staticcheck // SA1019 ignore deprecated GetOkExists
	// lintignore: XR001 // TODO: replace with alternative for GetOkExists
	if v, ok := d.GetOkExists("mentions_disabled"); ok {
		options.MentionsDisabled = gitlab.Bool(v.(bool))
	}

	if v, ok := d.GetOk("parent_id"); ok {
		options.ParentID = gitlab.Int(v.(int))
	}

	// nolint:staticcheck // SA1019 ignore deprecated GetOkExists
	// lintignore: XR001 // TODO: replace with alternative for GetOkExists
	if v, ok := d.GetOkExists("default_branch_protection"); ok {
		options.DefaultBranchProtection = gitlab.Int(v.(int))
	}

	log.Printf("[DEBUG] create gitlab group %q", *options.Name)

	group, _, err := client.Groups.CreateGroup(options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", group.ID))

	var updateOptions gitlab.UpdateGroupOptions

	// nolint:staticcheck // SA1019 ignore deprecated GetOkExists
	// lintignore: XR001 // TODO: replace with alternative for GetOkExists
	if v, ok := d.GetOkExists("prevent_forking_outside_group"); ok {
		updateOptions.PreventForkingOutsideGroup = gitlab.Bool(v.(bool))
	}

	if (updateOptions != gitlab.UpdateGroupOptions{}) {
		if _, _, err = client.Groups.UpdateGroup(d.Id(), &updateOptions, gitlab.WithContext(ctx)); err != nil {
			return diag.Errorf("could not update group after creation %q: %s", d.Id(), err)
		}
	}

	return resourceGitlabGroupRead(ctx, d, meta)
}

func resourceGitlabGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	log.Printf("[DEBUG] read gitlab group %s", d.Id())

	group, _, err := client.Groups.GetGroup(d.Id(), nil, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] gitlab group %s not found so removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	if group.MarkedForDeletionOn != nil {
		log.Printf("[DEBUG] gitlab group %s is marked for deletion", d.Id())
		d.SetId("")
		return nil
	}

	d.SetId(fmt.Sprintf("%d", group.ID))
	d.Set("name", group.Name)
	d.Set("path", group.Path)
	d.Set("full_path", group.FullPath)
	d.Set("full_name", group.FullName)
	d.Set("web_url", group.WebURL)
	d.Set("description", group.Description)
	d.Set("lfs_enabled", group.LFSEnabled)
	d.Set("request_access_enabled", group.RequestAccessEnabled)
	d.Set("visibility_level", group.Visibility)
	d.Set("project_creation_level", group.ProjectCreationLevel)
	d.Set("subgroup_creation_level", group.SubGroupCreationLevel)
	d.Set("require_two_factor_authentication", group.RequireTwoFactorAuth)
	d.Set("two_factor_grace_period", group.TwoFactorGracePeriod)
	d.Set("auto_devops_enabled", group.AutoDevopsEnabled)
	d.Set("emails_disabled", group.EmailsDisabled)
	d.Set("mentions_disabled", group.MentionsDisabled)
	d.Set("parent_id", group.ParentID)
	d.Set("runners_token", group.RunnersToken)
	d.Set("share_with_group_lock", group.ShareWithGroupLock)
	d.Set("default_branch_protection", group.DefaultBranchProtection)
	d.Set("prevent_forking_outside_group", group.PreventForkingOutsideGroup)

	return nil
}

func resourceGitlabGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	options := &gitlab.UpdateGroupOptions{}

	if d.HasChange("name") {
		options.Name = gitlab.String(d.Get("name").(string))
	}

	if d.HasChange("path") {
		options.Path = gitlab.String(d.Get("path").(string))
	}

	if d.HasChange("description") {
		options.Description = gitlab.String(d.Get("description").(string))
	}

	if d.HasChange("lfs_enabled") {
		options.LFSEnabled = gitlab.Bool(d.Get("lfs_enabled").(bool))
	}

	if d.HasChange("request_access_enabled") {
		options.RequestAccessEnabled = gitlab.Bool(d.Get("request_access_enabled").(bool))
	}

	// Always set visibility ; workaround for
	// https://gitlab.com/gitlab-org/gitlab-ce/issues/38459
	if v, ok := d.GetOk("visibility_level"); ok {
		options.Visibility = stringToVisibilityLevel(v.(string))
	}

	if d.HasChange("project_creation_level") {
		options.ProjectCreationLevel = stringToProjectCreationLevel(d.Get("project_creation_level").(string))
	}

	if d.HasChange("subgroup_creation_level") {
		options.SubGroupCreationLevel = stringToSubGroupCreationLevel(d.Get("subgroup_creation_level").(string))
	}

	if d.HasChange("require_two_factor_authentication") {
		options.RequireTwoFactorAuth = gitlab.Bool(d.Get("require_two_factor_authentication").(bool))
	}

	if d.HasChange("two_factor_grace_period") {
		options.TwoFactorGracePeriod = gitlab.Int(d.Get("two_factor_grace_period").(int))
	}

	if d.HasChange("auto_devops_enabled") {
		options.AutoDevopsEnabled = gitlab.Bool(d.Get("auto_devops_enabled").(bool))
	}

	if d.HasChange("emails_disabled") {
		options.EmailsDisabled = gitlab.Bool(d.Get("emails_disabled").(bool))
	}

	if d.HasChange("mentions_disabled") {
		options.MentionsDisabled = gitlab.Bool(d.Get("mentions_disabled").(bool))
	}

	if d.HasChange("share_with_group_lock") {
		options.ShareWithGroupLock = gitlab.Bool(d.Get("share_with_group_lock").(bool))
	}

	if d.HasChange("default_branch_protection") {
		options.DefaultBranchProtection = gitlab.Int(d.Get("default_branch_protection").(int))
	}

	if d.HasChange("prevent_forking_outside_group") {
		options.PreventForkingOutsideGroup = gitlab.Bool(d.Get("prevent_forking_outside_group").(bool))
	}

	log.Printf("[DEBUG] update gitlab group %s", d.Id())

	_, _, err := client.Groups.UpdateGroup(d.Id(), options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("parent_id") {
		err = transferSubGroup(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGitlabGroupRead(ctx, d, meta)
}

func transferSubGroup(ctx context.Context, d *schema.ResourceData, client *gitlab.Client) error {
	o, n := d.GetChange("parent_id")
	parentId, ok := n.(int)
	if !ok {
		return fmt.Errorf("error converting parent_id %v into an int", n)
	}

	opt := &gitlab.TransferSubGroupOptions{}
	if parentId != 0 {
		log.Printf("[DEBUG] transfer gitlab group %s from %v to new parent group %v", d.Id(), o, n)
		opt.GroupID = gitlab.Int(parentId)
	} else {
		log.Printf("[DEBUG] turn gitlab group %s from %v to a new top-level group", d.Id(), o)
	}

	_, _, err := client.Groups.TransferSubGroup(d.Id(), opt, gitlab.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("error transfering group %s to new parent group %v: %s", d.Id(), parentId, err)
	}

	return nil
}

func resourceGitlabGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	log.Printf("[DEBUG] Delete gitlab group %s", d.Id())

	_, err := client.Groups.DeleteGroup(d.Id(), gitlab.WithContext(ctx))
	if err != nil && !strings.Contains(err.Error(), "Group has been already marked for deletion") {
		return diag.Errorf("error deleting group %s: %s", d.Id(), err)
	}

	// Wait for the group to be deleted.
	// Deleting a group in gitlab is async.
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Deleting"},
		Target:  []string{"Deleted"},
		Refresh: func() (interface{}, string, error) {
			out, response, err := client.Groups.GetGroup(d.Id(), nil, gitlab.WithContext(ctx))
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					return out, "Deleted", nil
				}
				log.Printf("[ERROR] Received error: %#v", err)
				return out, "Error", err
			}
			if out.MarkedForDeletionOn != nil {
				// Represents a Gitlab EE soft-delete
				return out, "Deleted", nil
			}
			return out, "Deleting", nil
		},

		Timeout:    10 * time.Minute,
		MinTimeout: 3 * time.Second,
		Delay:      5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for group (%s) to become deleted: %s", d.Id(), err)
	}
	return nil
}
