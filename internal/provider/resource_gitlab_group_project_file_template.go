package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
	"log"
)

var _ = registerResource("gitlab_group_project_file_template", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_group_project_file_template`" + ` resource allows setting a project from which
custom file templates will be loaded. In order to use this resource, the project selected must be a direct child of
the group selected. After the resource has run, ` + "`gitlab_project_template.template_project_id`" + ` is available for use.
For more information about which file types are available as templates, view 
[GitLab's documentation](https://docs.gitlab.com/ee/user/group/custom_project_templates.html)

-> This resource requires a GitLab Enterprise instance with a Premium license.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/groups.html#update-group)`,

		// Since this resource updates an in-place resource, the update method is the same as the create method
		CreateContext: resourceGitLabGroupProjectFileTemplateCreateOrUpdate,
		UpdateContext: resourceGitLabGroupProjectFileTemplateCreateOrUpdate,
		ReadContext:   resourceGitLabGroupProjectFileTemplateRead,
		DeleteContext: resourceGitLabGroupProjectFileTemplateDelete,
		// Since this resource updates an in-place resource, importing doesn't make much sense. Simply add the resource
		// to the config and terraform will overwrite what's already in place and manage it from there.
		Schema: map[string]*schema.Schema{
			"group_id": {
				Description: `The ID of the group that will use the file template project. This group must be the direct
                parent of the project defined by project_id`,
				Type: schema.TypeInt,

				// Even though there is no traditional resource to create, leave "ForceNew" as "true" so that if someone
				// changes a configuration to a different group, the old group gets "deleted" (updated to have a value
				// of 0).
				ForceNew: true,
				Required: true,
			},
			"file_template_project_id": {
				Description: `The ID of the project that will be used for file templates. This project must be the direct
				child of the project defined by the group_id`,
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
})

func resourceGitLabGroupProjectFileTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	groupID := d.Get("group_id").(int)
	group, _, err := client.Groups.GetGroup(groupID, nil, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] gitlab group %d not found, removing from state", groupID)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	if group.MarkedForDeletionOn != nil {
		log.Printf("[DEBUG] gitlab group %s is marked for deletion, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.SetId(fmt.Sprintf("%d", group.ID))
	d.Set("file_template_project_id", group.FileTemplateProjectID)

	return nil
}

func resourceGitLabGroupProjectFileTemplateCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	groupID := d.Get("group_id").(int)
	projectID := gitlab.Int(d.Get("file_template_project_id").(int))

	// Creating the resource means updating the existing group to link the project to the group.
	options := &gitlab.UpdateGroupOptions{}
	if d.HasChanges("file_template_project_id") {
		options.FileTemplateProjectID = projectID
	}

	_, _, err := client.Groups.UpdateGroup(groupID, options)
	if err != nil {
		return diag.Errorf("unable to update group %d with `file_template_project_id` set to %d: %s", groupID, projectID, err)
	}
	return resourceGitLabGroupProjectFileTemplateRead(ctx, d, meta)
}

func resourceGitLabGroupProjectFileTemplateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	groupID := d.Get("group_id").(int)
	options := &gitlab.UpdateGroupOptions{}

	_, _, err := updateGroupWithOverwrittenFileTemplateOption(client, groupID, options)
	if err != nil {
		return diag.Errorf("could not update group %d to remove file template ID: %s", groupID, err)
	}
	return resourceGitLabGroupProjectFileTemplateRead(ctx, d, meta)
}

func updateGroupWithOverwrittenFileTemplateOption(client *gitlab.Client, groupID int, options *gitlab.UpdateGroupOptions) (*gitlab.Group, *gitlab.Response, error) {
	return client.Groups.UpdateGroup(groupID, options, func(request *retryablehttp.Request) error {
		//Overwrite the GroupUpdateOptions struct to remove the "omitempty", which forces the client to send an empty
		//string in just this request.
		removeOmitEmptyOptions := struct {
			FileTemplateProjectID *string `url:"file_template_project_id" json:"file_template_project_id"`
		}{
			FileTemplateProjectID: nil,
		}

		//Create the new body request with the above struct
		newBody, err := json.Marshal(removeOmitEmptyOptions)
		if err != nil {
			return err
		}

		//Set the request body to have the newly updated body
		err = request.SetBody(newBody)
		if err != nil {
			return err
		}

		return nil
	})
}
