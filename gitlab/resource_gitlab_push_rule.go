package gitlab

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabPushRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabPushRuleCreate,
		Read:   resourceGitlabPushRuleRead,
		Update: resourceGitlabPushRuleUpdate,
		Delete: resourceGitlabPushRuleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGitlabPushRuleStateImporter,
		},
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"deny_delete_tag": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"member_check": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"prevent_secrets": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"commit_message_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"branch_name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"author_email_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_file_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceGitlabPushRuleSetToState(d *schema.ResourceData, pushRules *gitlab.ProjectPushRules) {
	d.SetId(fmt.Sprintf("%d", pushRules.ID))
	d.Set("project", pushRules.ProjectID)
	d.Set("deny_delete_tag", pushRules.DenyDeleteTag)
	d.Set("member_check", pushRules.MemberCheck)
	d.Set("prevent_secrets", pushRules.PreventSecrets)
	d.Set("commit_message_regex", pushRules.CommitMessageRegex)
	d.Set("branch_name_regex", pushRules.BranchNameRegex)
	d.Set("author_email_regex", pushRules.AuthorEmailRegex)
	d.Set("file_name_regex", pushRules.FileNameRegex)
	d.Set("max_file_size", pushRules.MaxFileSize)
}

func resourceGitlabPushRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	options := &gitlab.AddProjectPushRuleOptions{}
	projectID := d.Get("project")

	if value, ok := d.GetOk("deny_delete_tag"); ok {
		options.DenyDeleteTag = gitlab.Bool(value.(bool))
	}

	if value, ok := d.GetOk("member_check"); ok {
		options.MemberCheck = gitlab.Bool(value.(bool))
	}

	if value, ok := d.GetOk("prevent_secrets"); ok {
		options.PreventSecrets = gitlab.Bool(value.(bool))
	}

	if value, ok := d.GetOk("commit_message_regex"); ok {
		options.CommitMessageRegex = gitlab.String(value.(string))
	}

	if value, ok := d.GetOk("branch_name_regex"); ok {
		options.BranchNameRegex = gitlab.String(value.(string))
	}

	if value, ok := d.GetOk("author_email_regex"); ok {
		options.AuthorEmailRegex = gitlab.String(value.(string))
	}

	if value, ok := d.GetOk("file_name_regex"); ok {
		options.FileNameRegex = gitlab.String(value.(string))
	}

	if value, ok := d.GetOk("max_file_size"); ok {
		options.MaxFileSize = gitlab.Int(value.(int))
	}

	log.Printf("[DEBUG] create gitlab push rule %q", projectID)

	pushRules, _, err := client.Projects.AddProjectPushRule(projectID, options)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", pushRules.ID))

	return resourceGitlabPushRuleRead(d, meta)
}

func resourceGitlabPushRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	log.Printf("[DEBUG] read gitlab push rule %s", d.Get("project"))

	pushRules, response, err := client.Projects.GetProjectPushRules(d.Get("project"))

	if err != nil {
		if response.StatusCode == 404 {
			log.Printf("[WARN] removing push rule %s from state because it no longer exists in gitlab", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	resourceGitlabPushRuleSetToState(d, pushRules)

	return nil
}

func resourceGitlabPushRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	options := &gitlab.EditProjectPushRuleOptions{}

	if d.HasChange("deny_delete_tag") {
		options.DenyDeleteTag = gitlab.Bool(d.Get("deny_delete_tag").(bool))
	}

	if d.HasChange("member_check") {
		options.MemberCheck = gitlab.Bool(d.Get("member_check").(bool))
	}

	if d.HasChange("prevent_secrets") {
		options.PreventSecrets = gitlab.Bool(d.Get("prevent_secrets").(bool))
	}

	if d.HasChange("commit_message_regex") {
		options.CommitMessageRegex = gitlab.String(d.Get("commit_message_regex").(string))
	}

	if d.HasChange("branch_name_regex") {
		options.BranchNameRegex = gitlab.String(d.Get("branch_name_regex").(string))
	}

	if d.HasChange("author_email_regex") {
		options.AuthorEmailRegex = gitlab.String(d.Get("author_email_regex").(string))
	}

	if d.HasChange("file_name_regex") {
		options.FileNameRegex = gitlab.String(d.Get("file_name_regex").(string))
	}

	if d.HasChange("max_file_size") {
		options.MaxFileSize = gitlab.Int(d.Get("max_file_size").(int))
	}

	if *options != (gitlab.EditProjectPushRuleOptions{}) {
		log.Printf("[DEBUG] update gitlab push rule %s", d.Get("project"))
		_, _, err := client.Projects.EditProjectPushRule(d.Get("project"), options)
		if err != nil {
			return err
		}
	}

	return resourceGitlabPushRuleRead(d, meta)
}

func resourceGitlabPushRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	log.Printf("[DEBUG] Delete gitlab push rule %s", d.Get("project"))

	_, err := client.Projects.DeleteProjectPushRule(d.Get("project"))
	if err != nil {
		return err
	}

	return nil
}

func resourceGitlabPushRuleStateImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	projectID := d.Id()

	client := meta.(*gitlab.Client)
	pushRules, _, err := client.Projects.GetProjectPushRules(projectID)

	if err != nil {
		return nil, err
	}

	d.SetId(fmt.Sprintf("%d", pushRules.ID))
	d.Set("project", projectID)

	return []*schema.ResourceData{d}, nil
}
