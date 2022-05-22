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

var managedLicenseAllowedValues = []string{
	"approved", "blacklisted", "allowed", "denied",
}

var _ = registerResource("gitlab_managed_license", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`" + `gitlab_managed_license` + "`" + ` resource allows to manage the lifecycle of a managed license.

-> This resource requires a GitLab Enterprise instance.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/managed_licenses.html)`,

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
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         false,
				ValidateFunc:     validation.StringInSlice(managedLicenseAllowedValues, true),
				DiffSuppressFunc: checkDeprecatedValuesForDiff,
				Description: fmt.Sprintf(`The approval status of the license. Valid values are: %s. "approved" and "blacklisted"
				have been deprecated in favor of "allowed" and "denied"; use "allowed" and "denied" for GitLab versions 15.0 and higher.
				Prior to version 15.0 and after 14.6, the values are equivalent.`, renderValueListForDocs(managedLicenseAllowedValues)),
			},
		},
	}
})

func resourceGitlabManagedLicenseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)

	approvalStatus, err := stringToApprovalStatus(ctx, client, d.Get("approval_status").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	options := &gitlab.AddManagedLicenseOptions{
		Name:           gitlab.String(d.Get("name").(string)),
		ApprovalStatus: approvalStatus,
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

	approvalStatus, err := stringToApprovalStatus(ctx, client, d.Get("approval_status").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	opts := &gitlab.EditManagedLicenceOptions{
		ApprovalStatus: approvalStatus,
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
func stringToApprovalStatus(ctx context.Context, client *gitlab.Client, s string) (*gitlab.LicenseApprovalStatusValue, error) {
	var value gitlab.LicenseApprovalStatusValue
	notSupported, err := isGitLabVersionAtLeast(ctx, client, "15.0")()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch GitLab version: %+v", err)
	}
	if notSupported {
		value = gitlab.LicenseApprovalStatusValue(s)
	} else {
		lookup := map[string]gitlab.LicenseApprovalStatusValue{
			"approved":    gitlab.LicenseApproved,
			"blacklisted": gitlab.LicenseBlacklisted,

			// This is counter-intuitive, but currently the API response from the non-deprecated
			// values is the deprecated values. So we have to map them here.
			"allowed": gitlab.LicenseApproved,
			"denied":  gitlab.LicenseBlacklisted,
		}

		v, ok := lookup[s]
		if !ok {
			return nil, fmt.Errorf("invalid approval status value %q", s)
		}
		value = v
	}
	return &value, nil
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

func checkDeprecatedValuesForDiff(k, oldValue, newValue string, d *schema.ResourceData) bool {
	approvedValues := []string{"approved", "allowed"}
	deniedValues := []string{"blacklisted", "denied"}

	// While we could technically combine these two "if" blocks, this seems more readable
	// and should have the same execution pattern.
	if contains(approvedValues, oldValue) && contains(approvedValues, newValue) {
		return true
	}

	if contains(deniedValues, oldValue) && contains(deniedValues, newValue) {
		return true
	}

	return false
}
