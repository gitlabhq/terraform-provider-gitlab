package provider

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/xanzy/go-gitlab"
)

func resourceGitlabManagedLicense() *schema.Resource {
	return &schema.Resource{
		Description: "This resource allows you to add rules for managing licenses on a project.\n" +
			"For additional information, please see the " +
			"[gitlab documentation](https://docs.gitlab.com/ee/user/compliance/license_compliance/).\n\n" +
			"~> Using this resource requires an active [gitlab ultimate](https://about.gitlab.com/pricing/)" +
			"subscription.",

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

	licenseId := strconv.Itoa(addManagedLicense.ID)
	d.SetId(buildTwoPartID(&project, &licenseId))

	return resourceGitlabManagedLicenseRead(ctx, d, meta)
}

func resourceGitlabManagedLicenseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project, licenseId, err := projectIdAndLicenseIdFromId(d.Id())
	if err != nil {
		return diag.FromErr(err)
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
	project, licenseId, err := projectIdAndLicenseIdFromId(d.Id())
	if err != nil {
		diag.FromErr(err)
	}

	log.Printf("[DEBUG] read gitlab Managed License for project/licenseId %s/%d", project, licenseId)
	license, _, err := client.ManagedLicenses.GetManagedLicense(project, licenseId, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] Managed License %s:%d no longer exists and is being removed from state", project, licenseId)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("project", project)
	d.Set("name", license.Name)
	d.Set("approval_status", license.ApprovalStatus)

	return nil
}

func resourceGitlabManagedLicenseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project, licenseId, err := projectIdAndLicenseIdFromId(d.Id())
	if err != nil {
		return diag.FromErr(err)
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

	licenseIdStr := strconv.Itoa(licenseId)
	d.SetId(buildTwoPartID(&project, &licenseIdStr))

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

func projectIdAndLicenseIdFromId(id string) (string, int, error) {
	projectId, id, err := parseTwoPartID(id)
	if err != nil {
		return "", 0, err
	}

	licenseId, err := strconv.Atoi(id)
	if err != nil {
		return "", 0, err
	}

	return projectId, licenseId, nil
}
