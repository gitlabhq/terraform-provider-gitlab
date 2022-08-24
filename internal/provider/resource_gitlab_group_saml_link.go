package provider

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_group_saml_link", func() *schema.Resource {
	validGroupSamlLinkAccessLevelNames := []string{
		"guest",
		"reporter",
		"developer",
		"maintainer",
		"owner",
	}

	return &schema.Resource{
		Description: `The ` + "`gitlab_group_saml_link`" + ` resource allows to manage the lifecycle of an SAML integration with a group.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/groups.html#saml-group-links)`,

		CreateContext: resourceGitlabGroupSamlLinkCreate,
		ReadContext:   resourceGitlabGroupSamlLinkRead,
		DeleteContext: resourceGitlabGroupSamlLinkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"group": {
				Description: "The ID or path of the group to add the SAML Group Link to.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"saml_group_name": {
				Description: "The name of the SAML group.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"access_level": {
				Description:      fmt.Sprintf("Access level for members of the SAML group. Valid values are: %s.", renderValueListForDocs(validGroupSamlLinkAccessLevelNames)),
				Type:             schema.TypeString,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(validGroupSamlLinkAccessLevelNames, false)),
				Required:         true,
				ForceNew:         true,
			},
		},
	}
})

func resourceGitlabGroupSamlLinkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	group := d.Get("group").(string)
	samlGroupName := d.Get("saml_group_name").(string)
	accessLevel := accessLevelNameToValue[d.Get("access_level").(string)]

	options := &gitlab.AddGroupSAMLLinkOptions{
		SAMLGroupName: gitlab.String(samlGroupName),
		AccessLevel:   gitlab.AccessLevel(accessLevel),
	}

	log.Printf("[DEBUG] Create GitLab Group SAML Link for group %q with name %q", group, samlGroupName)
	SamlLink, _, err := client.Groups.AddGroupSAMLLink(group, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(buildTwoPartID(&group, &SamlLink.Name))
	return resourceGitlabGroupSamlLinkRead(ctx, d, meta)
}

func resourceGitlabGroupSamlLinkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	group, samlGroupName, parse_err := parseTwoPartID(d.Id())
	if parse_err != nil {
		return diag.FromErr(parse_err)
	}

	// Try to fetch all group links from GitLab
	log.Printf("[DEBUG] Read GitLab Group SAML Link for group %q", group)
	samlLink, _, err := client.Groups.GetGroupSAMLLink(group, samlGroupName, nil, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] GitLab SAML Group Link %s for group ID %s not found, removing from state", samlGroupName, group)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("group", group)
	d.Set("access_level", accessLevelValueToName[samlLink.AccessLevel])
	d.Set("saml_group_name", samlLink.Name)

	return nil
}

func resourceGitlabGroupSamlLinkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	group, samlGroupName, parse_err := parseTwoPartID(d.Id())
	if parse_err != nil {
		return diag.FromErr(parse_err)
	}

	log.Printf("[DEBUG] Delete GitLab Group SAML Link for group %q with name %q", group, samlGroupName)
	_, err := client.Groups.DeleteGroupSAMLLink(group, samlGroupName, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[WARNING] %s", err)
		} else {
			return diag.FromErr(err)
		}
	}

	return nil
}
