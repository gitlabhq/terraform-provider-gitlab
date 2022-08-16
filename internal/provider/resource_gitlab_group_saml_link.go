package provider

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_group_saml_link", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_group_saml_link`" + ` resource allows to manage the lifecycle of an SAML integration with a group.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/groups.html#saml-group-links)`,

		CreateContext: resourceGitlabGroupSamlLinkCreate,
		ReadContext:   resourceGitlabGroupSamlLinkRead,
		DeleteContext: resourceGitlabGroupSamlLinkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGitlabGroupSamlLinkImporter,
		},

		Schema: map[string]*schema.Schema{
			"group_id": {
				Description: "The id of the GitLab group.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"access_level": {
				Description:      fmt.Sprintf("Minimum access level for members of the SAML group. Valid values are: %s", renderValueListForDocs(validGroupSamlLinkAccessLevelNames)),
				Type:             schema.TypeString,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(validGroupSamlLinkAccessLevelNames, false)),
				Required:         true,
				ForceNew:         true,
			},
			"saml_group_name": {
				Description: "The name of the SAML group.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"force": {
				Description: "If true, then delete and replace an existing SAML link if one exists.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
			},
		},
	}
})

func resourceGitlabGroupSamlLinkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	groupId := d.Get("group_id").(string)
	accessLevel := d.Get("access_level").(string)
	samlGroupName := d.Get("saml_group_name").(string)
	force := d.Get("force").(bool)

	options := &gitlab.AddGroupSAMLLinkOptions{
		AccessLevel:   &accessLevel,
		SamlGroupName: &samlGroupName,
	}

	if force {
		if err := resourceGitlabGroupSamlLinkDelete(ctx, d, meta); err != nil {
			return err
		}
	}

	log.Printf("[DEBUG] Create GitLab group SamlLink %s", d.Id())
	SamlLink, _, err := client.Groups.AddGroupSAMLLink(groupId, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(buildTwoPartID(&groupId, &SamlLink.Name))

	return resourceGitlabGroupSamlLinkRead(ctx, d, meta)
}

func resourceGitlabGroupSamlLinkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	groupId := d.Get("group_id").(string)

	// Try to fetch all group links from GitLab
	log.Printf("[DEBUG] Read GitLab group SamlLinks %s", groupId)
	samlLinks, _, err := client.Groups.ListGroupSAMLLinks(groupId, nil, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	// If we got here and don't have links, assume GitLab is below version 12.8 and skip the check
	if samlLinks != nil {
		// Check if the LDAP link exists in the returned list of links
		found := false
		for _, samlLink := range samlLinks {
			if buildTwoPartID(&groupId, &samlLink.Name) == d.Id() {
				d.Set("group_id", groupId)
				d.Set("access_level", samlLink.AccessLevel)
				d.Set("saml_group_name", samlLink.Name)
				found = true
				break
			}
		}

		if !found {
			d.SetId("")
			return diag.Errorf("SamlLink %s does not exist.", d.Id())
		}
	}

	return nil
}

func resourceGitlabGroupSamlLinkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	groupId := d.Get("group_id").(string)
	samlGroupName := d.Get("saml_group_name").(string)

	log.Printf("[DEBUG] Delete GitLab group SamlLink %s", d.Id())
	_, err := client.Groups.DeleteGroupSAMLLink(groupId, samlGroupName, cn, gitlab.WithContext(ctx))
	if err != nil {
		switch err.(type) { // nolint // TODO: Resolve this golangci-lint issue: S1034: assigning the result of this type assertion to a variable (switch err := err.(type)) could eliminate type assertions in switch cases (gosimple)
		case *gitlab.ErrorResponse:
			// Ignore SAML links that don't exist
			if strings.Contains(string(err.(*gitlab.ErrorResponse).Message), "Linked SAML group not found") { // nolint // TODO: Resolve this golangci-lint issue: S1034(related information): could eliminate this type assertion (gosimple)
				log.Printf("[WARNING] %s", err)
			} else {
				return diag.FromErr(err)
			}
		default:
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceGitlabGroupSamlLinkImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid saml link import id (should be <group id>:<saml group name>): %s", d.Id())
	}

	groupId, samlGroupName := parts[0], parts[1]
	d.SetId(buildTwoPartID(&groupId, &samlGroupName))
	d.Set("group_id", groupId)
	d.Set("force", false)

	diag := resourceGitlabGroupSamlLinkRead(ctx, d, meta)
	if diag.HasError() {
		return nil, fmt.Errorf("%s", diag[0].Summary)
	}
	return []*schema.ResourceData{d}, nil
}
