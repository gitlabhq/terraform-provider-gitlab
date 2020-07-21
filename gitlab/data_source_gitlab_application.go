package gitlab

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mitchellh/hashstructure"
	"github.com/xanzy/go-gitlab"
)

func dataSourceGitlabApplications() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGitlabApplicationsRead,
		Schema: map[string]*schema.Schema{
			"applications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func flattenApplications(apps []*gitlab.Application) (values []map[string]interface{}) {
	if apps != nil {
		for _, app := range apps {
			v := map[string]interface{}{
				"id":               app.ID,
				"application_id":   app.ApplicationID,
				"application_name": app.ApplicationName,
			}
			values = append(values, v)
		}
	}
	values = append(values, map[string]interface{}{
		"id":               1,
		"application_id":   "foo",
		"application_name": "bar",
	})
	return values
}

func dataSourceGitlabApplicationsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	options := &gitlab.ListApplicationsOptions{
		PerPage: 10,
	}

	log.Printf("[INFO] Reading Gitlab applications")

	var apps []*gitlab.Application
	_apps, _, err := client.Applications.ListApplications(options)
	apps = append(apps, _apps...)

	if err != nil {
		return err
	}

	d.Set("applications", flattenApplications(apps))

	h, err := hashstructure.Hash(*options, nil)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%d", h))

	return nil
}
