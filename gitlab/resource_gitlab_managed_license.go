package gitlab

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/xanzy/go-gitlab"
	"log"
	"strconv"
)

func resourceGitlabManagedLicense() *schema.Resource {
	return &schema.Resource{
		Description: "This resource allows you to add rules for managing licenses on a project.\n" +
			"For additional information, please see the " +
			"[gitlab documentation](https://docs.gitlab.com/ee/user/compliance/license_compliance/).",

		CreateContext: resourceGitlabManagedLicenseCreate,
		ReadContext:   resourceGitlabManagedLicenseRead,
		UpdateContext: resourceGitlabManagedLicenseUpdate,
		DeleteContext: resourceGitlabManagedLicenseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the project under which the managed license will be created.",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				// GitLab's edit API doesn't allow us to edit the name, only
				// the approval status.
				ForceNew:    true,
				Description: "The name of the managed license (I.e., 'Apache License 2.0' or 'MIT license')",
			},
			"approval_status": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validation.StringInSlice([]string{"approved", "blacklisted"}, true),
				Description:  "Whether the license is approved or not. Only 'approved' or 'blacklisted' allowed.",
			},
		},
	}
}

func resourceGitlabManagedLicenseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)

	options := &gitlab.AddManagedLicenseOptions{
		Name:           gitlab.String(d.Get("name").(string)),
		ApprovalStatus: stringToApprovalStatus(d.Get("approval_status").(string)),
	}

	log.Printf("[DEBUG] create gitlab Managed License on Project %s, with Name %s", project, *options.Name)

	addManagedLicense, _, err := client.ManagedLicenses.AddManagedLicense(project, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(addManagedLicense.ID))
	return resourceGitlabManagedLicenseRead(ctx, d, meta)
}

func resourceGitlabManagedLicenseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	licenseId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s cannot be converted to int", d.Id()))
	}

	log.Printf("[DEBUG] Delete gitlab Managed License %s", d.Id())
	_, err = client.ManagedLicenses.DeleteManagedLicense(project, licenseId, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete managed license %d for projec %s: %w", licenseId, project, err))
	}

	return nil
}

func resourceGitlabManagedLicenseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	id, err := strconv.Atoi(d.Id())

	if err != nil {
		return diag.FromErr(fmt.Errorf("%s cannot be converted to int", d.Id()))
	}
	log.Printf("[DEBUG] read gitlab Managed License for project/id %s/%d", project, id)

	license, _, err := client.ManagedLicenses.GetManagedLicense(project, id, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s cannot be converted to int", d.Id()))
	}

	d.Set("project", license.ID)
	d.Set("name", license.Name)
	d.Set("approval_status", license.ApprovalStatus)

	return nil
}

func resourceGitlabManagedLicenseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	licenseId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s cannot be converted to int", d.Id()))
	}

	opts := &gitlab.EditManagedLicenceOptions{
		ApprovalStatus: stringToApprovalStatus(d.Get("approval_status").(string)),
	}

	if d.HasChange("approval_status") {
		opts.ApprovalStatus = stringToApprovalStatus(d.Get("approval_status").(string))
	}

	log.Printf("[DEBUG] update gitlab Managed License %s", d.Id())
	_, _, err = client.ManagedLicenses.EditManagedLicense(project, licenseId, opts, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGitlabManagedLicenseRead(ctx, d, meta)
}

// Convert the incoming string into the proper constant value for passing into the API.
func stringToApprovalStatus(s string) *gitlab.LicenseApprovalStatusValue {
	lookup := map[string]gitlab.LicenseApprovalStatusValue{
		"approved":    gitlab.LicenseApproved,
		"blacklisted": gitlab.LicenseBlacklisted,
	}

	value, ok := lookup[s]
	if !ok {
		return nil
	}
	return &value
}
