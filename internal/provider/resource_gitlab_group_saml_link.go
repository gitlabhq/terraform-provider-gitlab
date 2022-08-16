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

var validGroupSamlLinkAccessLevelNames = []string{
	"Guest",
	"Reporter",
	"Developer",
	"Maintainer",
	"Owner",
}

var _ = registerResource("gitlab_group_saml_link", func() *schema.Resource {
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
		},
	}
})

func resourceGitlabGroupSamlLinkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	group := d.Get("group").(string)
	accessLevel := d.Get("access_level").(string)
	samlGroupName := d.Get("saml_group_name").(string)

	options := &gitlab.AddGroupSAMLLinkOptions{
		AccessLevel:   gitlab.String(accessLevel),
		SamlGroupName: gitlab.String(samlGroupName),
	}

	log.Printf("[DEBUG] Create GitLab group SamlLink %s", d.Id())
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
	log.Printf("[DEBUG] Read GitLab group SamlLinks %s", group)
	samlLinks, _, err := client.Groups.ListGroupSAMLLinks(group, nil, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	if samlLinks != nil {
		// Check if the SAML link exists in the returned list of links
		found := false
		for _, samlLink := range samlLinks {
			if samlLink.Name == samlGroupName {
				d.Set("group", group)
				d.Set("access_level", samlLink.AccessLevel)
				d.Set("saml_group_name", samlLink.Name)
				found = true
				break
			}
		}

		if !found {
			log.Printf("[DEBUG] GitLab SAML Group Link %d, group ID %s not found, removing from state", samlGroupName, group)
			d.SetId("")
			return nil
		}
	}

	return nil
}

func resourceGitlabGroupSamlLinkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	group, samlGroupName, parse_err := parseTwoPartID(d.Id())
	if parse_err != nil {
		return diag.FromErr(parse_err)
	}

	log.Printf("[DEBUG] Delete GitLab group SamlLink %s", d.Id())
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
