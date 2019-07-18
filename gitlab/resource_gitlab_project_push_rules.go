package gitlab

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabProjectPushRules() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabProjectPushRulesCreate,
		Read:   resourceGitlabProjectPushRulesRead,
		Update: resourceGitlabProjectPushRulesUpdate,
		Delete: resourceGitlabProjectPushRulesDelete,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"commit_message_regex": {
				Type:     schema.TypeString,
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
