package gitlab

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabProjectFreezePeriod() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabProjectFreezePeriodCreate,
		Read:   resourceGitlabProjectFreezePeriodRead,
		Update: resourceGitlabProjectFreezePeriodUpdate,
		Delete: resourceGitlabProjectFreezePeriodDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"freeze_start": {
				Type:     schema.TypeString,
				Required: true,
			},
			"freeze_end": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cron_timezone": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "UTC",
			},
		},
	}
}

func resourceGitlabProjectFreezePeriodCreate(d *schema.ResourceData, meta interface{}) error {
	projectID := d.Get("project_id").(string)

	options := gitlab.CreateFreezePeriodOptions{
		FreezeStart:  gitlab.String(d.Get("freeze_start").(string)),
		FreezeEnd:    gitlab.String(d.Get("freeze_end").(string)),
		CronTimezone: gitlab.String(d.Get("cron_timezone").(string)),
	}

	log.Printf("[DEBUG] Project %s create gitlab project-level freeze period %+v", projectID, options)

	client := meta.(*gitlab.Client)
	FreezePeriod, _, err := client.FreezePeriods.CreateFreezePeriodOptions(projectID, &options)
	if err != nil {
		return err
	}

	FreezePeriodIDString := fmt.Sprintf("%d", FreezePeriod.ID)
	d.SetId(buildTwoPartID(&projectID, &FreezePeriodIDString))

	return resourceGitlabProjectFreezePeriodRead(d, meta)
}

func resourceGitlabProjectFreezePeriodRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	projectID, freezePeriodID, err := projectIDAndFreezePeriodIDFromID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] read gitlab FreezePeriod %s/%d", projectID, freezePeriodID)

	freezePeriod, resp, err := client.FreezePeriods.GetFreezePeriod(projectID, freezePeriodID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			log.Printf("[DEBUG] project freeze period for %s not found so removing it from state", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("freeze_start", freezePeriod.FreezeStart)
	d.Set("freeze_end", freezePeriod.FreezeEnd)
	d.Set("cron_timezone", freezePeriod.CronTimezone)
	d.Set("project_id", projectID)

	return nil
}

func resourceGitlabProjectFreezePeriodUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	projectID, freezePeriodID, err := projectIDAndFreezePeriodIDFromID(d.Id())
	options := &gitlab.UpdateFreezePeriodOptions{}

	if err != nil {
		return fmt.Errorf("%s cannot be converted to int", d.Id())
	}

	if d.HasChange("freeze_start") {
		options.FreezeStart = gitlab.String(d.Get("freeze_start").(string))
	}

	if d.HasChange("freeze_end") {
		options.FreezeEnd = gitlab.String(d.Get("freeze_end").(string))
	}

	if d.HasChange("cron_timezone") {
		options.CronTimezone = gitlab.String(d.Get("cron_timezone").(string))
	}

	log.Printf("[DEBUG] update gitlab FreezePeriod %s", d.Id())

	_, _, err = client.FreezePeriods.UpdateFreezePeriodOptions(projectID, freezePeriodID, options)
	if err != nil {
		return err
	}

	return resourceGitlabProjectFreezePeriodRead(d, meta)
}

func resourceGitlabProjectFreezePeriodDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	projectID, freezePeriodID, err := projectIDAndFreezePeriodIDFromID(d.Id())
	log.Printf("[DEBUG] Delete gitlab FreezePeriod %s", d.Id())

	if err != nil {
		return fmt.Errorf("%s cannot be converted to int", d.Id())
	}

	if _, err = client.FreezePeriods.DeleteFreezePeriod(projectID, freezePeriodID); err != nil {
		return fmt.Errorf("failed to delete pipeline schedule %q: %w", d.Id(), err)
	}

	return nil
}

func projectIDAndFreezePeriodIDFromID(id string) (string, int, error) {
	project, freezePeriodIDString, err := parseTwoPartID(id)
	if err != nil {
		return "", 0, err
	}

	freezePeriodID, err := strconv.Atoi(freezePeriodIDString)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get freezePeriodId: %v", err)
	}

	return project, freezePeriodID, nil
}
