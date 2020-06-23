package gitlab

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabProjectPushRules() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabProjectPushRulesCreate,
		Read:   resourceGitlabProjectPushRulesRead,
		Update: resourceGitlabProjectPushRulesUpdate,
		Delete: resourceGitlabProjectPushRulesDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGitlabProjectPushRulesImport,
		},
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"commit_message_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			/* Not implemented in gitlab client
			"commit_message_negative_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			*/
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
			"max_file_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}
func resourceGitlabProjectPushRulesUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	options := &gitlab.EditProjectPushRuleOptions{
		CommitMessageRegex: gitlab.String(d.Get("commit_message_regex").(string)),
		BranchNameRegex:    gitlab.String(d.Get("branch_name_regex").(string)),
		AuthorEmailRegex:   gitlab.String(d.Get("author_email_regex").(string)),
		FileNameRegex:      gitlab.String(d.Get("file_name_regex").(string)),
		DenyDeleteTag:      gitlab.Bool(d.Get("deny_delete_tag").(bool)),
		MemberCheck:        gitlab.Bool(d.Get("member_check").(bool)),
		PreventSecrets:     gitlab.Bool(d.Get("prevent_secrets").(bool)),
		MaxFileSize:        gitlab.Int(d.Get("max_file_size").(int)),
	}
	log.Printf("[DEBUG] update gitlab project %s push rules %#v", project, *options)
	_, _, err := client.Projects.EditProjectPushRule(project, options)
	if err != nil {
		return err
	}
	return resourceGitlabProjectPushRulesRead(d, meta)
}

func resourceGitlabProjectPushRulesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	options := &gitlab.AddProjectPushRuleOptions{
		CommitMessageRegex: gitlab.String(d.Get("commit_message_regex").(string)),
		BranchNameRegex:    gitlab.String(d.Get("branch_name_regex").(string)),
		AuthorEmailRegex:   gitlab.String(d.Get("author_email_regex").(string)),
		FileNameRegex:      gitlab.String(d.Get("file_name_regex").(string)),
		DenyDeleteTag:      gitlab.Bool(d.Get("deny_delete_tag").(bool)),
		MemberCheck:        gitlab.Bool(d.Get("member_check").(bool)),
		PreventSecrets:     gitlab.Bool(d.Get("prevent_secrets").(bool)),
		MaxFileSize:        gitlab.Int(d.Get("max_file_size").(int)),
	}
	log.Printf("[DEBUG] create gitlab project %s push rules %#v", project, *options)

	pushRules, _, err := client.Projects.AddProjectPushRule(project, options)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%d", pushRules.ID))
	return resourceGitlabProjectPushRulesRead(d, meta)
}

func resourceGitlabProjectPushRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	log.Printf("[DEBUG] read gitlab project %s", project)
	pushRules, _, err := client.Projects.GetProjectPushRules(project)
	if err != nil {
		return err
	}
	d.Set("commit_message_regex", pushRules.CommitMessageRegex)
	d.Set("branch_name_regex", pushRules.BranchNameRegex)
	d.Set("author_email_regex", pushRules.AuthorEmailRegex)
	d.Set("file_name_regex", pushRules.FileNameRegex)
	d.Set("deny_delete_tag", pushRules.DenyDeleteTag)
	d.Set("member_check", pushRules.MemberCheck)
	d.Set("prevent_secrets", pushRules.PreventSecrets)
	d.Set("max_file_size", pushRules.MaxFileSize)
	return nil
}

func resourceGitlabProjectPushRulesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	log.Printf("[DEBUG] Delete gitlab project push rules %s", project)
	log.Println(project)
	_, err := client.Projects.DeleteProjectPushRule(project)
	return err
}

func resourceGitlabProjectPushRulesImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.Set("project", d.Id())
	err := resourceGitlabProjectPushRulesRead(d, meta)
	return []*schema.ResourceData{d}, err
}
