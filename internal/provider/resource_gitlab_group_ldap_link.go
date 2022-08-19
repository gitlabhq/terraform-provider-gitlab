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

var _ = registerResource("gitlab_group_ldap_link", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_group_ldap_link`" + ` resource allows to manage the lifecycle of an LDAP integration with a group.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/groups.html#ldap-group-links)`,

		CreateContext: resourceGitlabGroupLdapLinkCreate,
		ReadContext:   resourceGitlabGroupLdapLinkRead,
		DeleteContext: resourceGitlabGroupLdapLinkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGitlabGroupLdapLinkImporter,
		},

		Schema: map[string]*schema.Schema{
			"group_id": {
				Description: "The id of the GitLab group.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"cn": {
				Description: "The CN of the LDAP group to link with.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"access_level": {
				Description:      fmt.Sprintf("Minimum access level for members of the LDAP group. Valid values are: %s", renderValueListForDocs(validGroupAccessLevelNames)),
				Type:             schema.TypeString,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(validGroupAccessLevelNames, false)),
				Optional:         true,
				ForceNew:         true,
				Deprecated:       "Use `group_access` instead of the `access_level` attribute.",
				ExactlyOneOf:     []string{"access_level", "group_access"},
			},
			"group_access": {
				Description:      fmt.Sprintf("Minimum access level for members of the LDAP group. Valid values are: %s", renderValueListForDocs(validGroupAccessLevelNames)),
				Type:             schema.TypeString,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(validGroupAccessLevelNames, false)),
				Optional:         true,
				ForceNew:         true,
				ExactlyOneOf:     []string{"access_level", "group_access"},
			},
			// Changing GitLab API parameter "provider" to "ldap_provider" to avoid clashing with the Terraform "provider" key word
			"ldap_provider": {
				Description: "The name of the LDAP provider as stored in the GitLab database. Note that this is NOT the value of the `label` attribute as shown in the web UI. In most cases this will be `ldapmain` but you may use the [LDAP check rake task](https://docs.gitlab.com/ee/administration/raketasks/ldap.html#check) for receiving the LDAP server name: `LDAP: ... Server: ldapmain`",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"force": {
				Description: "If true, then delete and replace an existing LDAP link if one exists.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
			},
		},
	}
})

func resourceGitlabGroupLdapLinkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	groupId := d.Get("group_id").(string)
	cn := d.Get("cn").(string)

	var groupAccess gitlab.AccessLevelValue
	if v, ok := d.GetOk("group_access"); ok {
		groupAccess = gitlab.AccessLevelValue(accessLevelNameToValue[v.(string)])
	} else if v, ok := d.GetOk("access_level"); ok {
		groupAccess = gitlab.AccessLevelValue(accessLevelNameToValue[v.(string)])
	} else {
		return diag.Errorf("Neither `group_access` nor `access_level` (deprecated) is set")
	}

	ldap_provider := d.Get("ldap_provider").(string)
	force := d.Get("force").(bool)

	options := &gitlab.AddGroupLDAPLinkOptions{
		CN:          &cn,
		GroupAccess: &groupAccess,
		Provider:    &ldap_provider,
	}

	if force {
		if err := resourceGitlabGroupLdapLinkDelete(ctx, d, meta); err != nil {
			return err
		}
	}

	log.Printf("[DEBUG] Create GitLab group LdapLink %s", d.Id())
	LdapLink, _, err := client.Groups.AddGroupLDAPLink(groupId, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(buildTwoPartID(&LdapLink.Provider, &LdapLink.CN))

	return resourceGitlabGroupLdapLinkRead(ctx, d, meta)
}

func resourceGitlabGroupLdapLinkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	groupId := d.Get("group_id").(string)

	// Try to fetch all group links from GitLab
	log.Printf("[DEBUG] Read GitLab group LdapLinks %s", groupId)
	ldapLinks, _, err := client.Groups.ListGroupLDAPLinks(groupId, nil, gitlab.WithContext(ctx))
	if err != nil {
		// The read/GET API wasn't implemented in GitLab until version 12.8 (March 2020, well after the add and delete APIs).
		// If we 404, assume GitLab is at an older version and take things on faith.
		switch err.(type) { // nolint // TODO: Resolve this golangci-lint issue: S1034: assigning the result of this type assertion to a variable (switch err := err.(type)) could eliminate type assertions in switch cases (gosimple)
		case *gitlab.ErrorResponse:
			if err.(*gitlab.ErrorResponse).Response.StatusCode == 404 { // nolint // TODO: Resolve this golangci-lint issue: S1034(related information): could eliminate this type assertion (gosimple)
				log.Printf("[WARNING] This GitLab instance doesn't have the GET API for group_ldap_sync.  Please upgrade to 12.8 or later for best results.")
			} else {
				return diag.FromErr(err)
			}
		default:
			return diag.FromErr(err)
		}
	}

	// If we got here and don't have links, assume GitLab is below version 12.8 and skip the check
	if ldapLinks != nil {
		// Check if the LDAP link exists in the returned list of links
		found := false
		for _, ldapLink := range ldapLinks {
			if buildTwoPartID(&ldapLink.Provider, &ldapLink.CN) == d.Id() {
				d.Set("group_id", groupId)
				d.Set("cn", ldapLink.CN)
				d.Set("group_access", accessLevelValueToName[ldapLink.GroupAccess])
				d.Set("ldap_provider", ldapLink.Provider)
				found = true
				break
			}
		}

		if !found {
			d.SetId("")
			return diag.Errorf("LdapLink %s does not exist.", d.Id())
		}
	}

	return nil
}

func resourceGitlabGroupLdapLinkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	groupId := d.Get("group_id").(string)
	cn := d.Get("cn").(string)
	ldap_provider := d.Get("ldap_provider").(string)

	log.Printf("[DEBUG] Delete GitLab group LdapLink %s", d.Id())
	_, err := client.Groups.DeleteGroupLDAPLinkForProvider(groupId, ldap_provider, cn, gitlab.WithContext(ctx))
	if err != nil {
		switch err.(type) { // nolint // TODO: Resolve this golangci-lint issue: S1034: assigning the result of this type assertion to a variable (switch err := err.(type)) could eliminate type assertions in switch cases (gosimple)
		case *gitlab.ErrorResponse:
			// Ignore LDAP links that don't exist
			if strings.Contains(string(err.(*gitlab.ErrorResponse).Message), "Linked LDAP group not found") { // nolint // TODO: Resolve this golangci-lint issue: S1034(related information): could eliminate this type assertion (gosimple)
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

func resourceGitlabGroupLdapLinkImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), ":", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid ldap link import id (should be <group id>:<ldap provider>:<ladp cn>): %s", d.Id())
	}

	groupId, ldapProvider, ldapCN := parts[0], parts[1], parts[2]
	d.SetId(buildTwoPartID(&ldapProvider, &ldapCN))
	d.Set("group_id", groupId)
	d.Set("force", false)

	diag := resourceGitlabGroupLdapLinkRead(ctx, d, meta)
	if diag.HasError() {
		return nil, fmt.Errorf("%s", diag[0].Summary)
	}
	return []*schema.ResourceData{d}, nil
}
