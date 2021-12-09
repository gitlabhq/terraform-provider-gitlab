package gitlab

import (
	"log"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func resourceGitlabTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabTopicCreate,
		Read:   resourceGitlabTopicRead,
		Update: resourceGitlabTopicUpdate,
		Delete: resourceGitlabTopicDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceGitlabTopicCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	options := &gitlab.CreateTopicOptions{
		Name: gitlab.String(d.Get("name").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		options.Description = gitlab.String(v.(string))
	}

	log.Printf("[DEBUG] create gitlab topic %s", *options.Name)

	topic, _, err := client.Topics.CreateTopic(options)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(int64(topic.ID), 10))

	return resourceGitlabTopicRead(d, meta)
}

func resourceGitlabTopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	topicID := d.Id()
	log.Printf("[DEBUG] read gitlab topic %s", topicID)

	topic, resp, err := client.Topics.GetTopic(topicID, nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			log.Printf("[DEBUG] gitlab group %s not found so removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(strconv.FormatInt(int64(topic.ID), 10))
	d.Set("name", topic.Name)
	d.Set("description", topic.Description)

	return nil
}

func resourceGitlabTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	options := &gitlab.UpdateTopicOptions{}

	if d.HasChange("name") {
		options.Name = gitlab.String(d.Get("name").(string))
	}

	if d.HasChange("description") {
		options.Description = gitlab.String(d.Get("description").(string))
	}

	log.Printf("[DEBUG] update gitlab topic %s", d.Id())

	_, _, err := client.Topics.UpdateTopic(d.Id(), options)
	if err != nil {
		return err
	}

	return resourceGitlabTopicRead(d, meta)
}

func resourceGitlabTopicDelete(d *schema.ResourceData, meta interface{}) error {

	// TODO: Use outcommented code below as soon as deleting a topic is supported
	log.Printf("[WARN] Not deleting gitlab topic %s as gitlab API doens't support deleting topics", d.Id())
	return nil
	/*
		client := meta.(*gitlab.Client)
		log.Printf("[DEBUG] Delete gitlab topic %s", d.Id())

		_, err := client.Topics.DeleteTopic(d.Id())
		return err

	*/
}
