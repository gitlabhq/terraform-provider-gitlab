package gitlab

import (
	"log"
	"net/url"

	"github.com/hashicorp/terraform/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabTagProtection() *schema.Resource {
	acceptedAccessLevels := make([]string, 0, len(accessLevelID))

	for k := range accessLevelID {
		acceptedAccessLevels = append(acceptedAccessLevels, k)
	}
	return &schema.Resource{
		Create: resourceGitlabTagProtectionCreate,
		Read:   resourceGitlabTagProtectionRead,
		Delete: resourceGitlabTagProtectionDelete,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tag": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"create_access_level": {
				Type:         schema.TypeString,
				ValidateFunc: validateValueFunc(acceptedAccessLevels),
				Required:     true,
				ForceNew:     true,
			},
		},
	}
}

func resourceGitlabTagProtectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	tag := gitlab.String(d.Get("tag").(string))
	createAccessLevel := accessLevelID[d.Get("create_access_level").(string)]

	options := &gitlab.ProtectRepositoryTagsOptions{
		Name:              tag,
		CreateAccessLevel: &createAccessLevel,
	}

	log.Printf("[DEBUG] create gitlab tag protection on %v for project %s", options.Name, project)

	tp, _, err := client.ProtectedTags.ProtectRepositoryTags(project, options)
	if err != nil {
		// Remove existing tag protection
		_, err = client.ProtectedTags.UnprotectRepositoryTags(project, url.PathEscape(*tag))
		if err != nil {
			return err
		}
		// Reprotect tag with updated values
		tp, _, err = client.ProtectedTags.ProtectRepositoryTags(project, options)
		if err != nil {
			return err
		}
	}

	d.SetId(buildTwoPartID(&project, &tp.Name))

	return resourceGitlabTagProtectionRead(d, meta)
}

func resourceGitlabTagProtectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project, tag, err := projectAndTagFromID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] read gitlab tag protection for project %s, tag %s", project, tag)

	pt, _, err := client.ProtectedTags.GetProtectedTag(project, url.PathEscape(tag))
	if err != nil {
		return err
	}

	d.Set("project", project)
	d.Set("tag", pt.Name)
	d.Set("create_access_level", pt.CreateAccessLevels[0].AccessLevel)

	d.SetId(buildTwoPartID(&project, &pt.Name))

	return nil
}

func resourceGitlabTagProtectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	tag := d.Get("tag").(string)

	log.Printf("[DEBUG] Delete gitlab protected tag %s for project %s", tag, project)

	_, err := client.ProtectedTags.UnprotectRepositoryTags(project, url.PathEscape(tag))
	return err
}

func projectAndTagFromID(id string) (string, string, error) {
	project, tag, err := parseTwoPartID(id)

	if err != nil {
		log.Printf("[WARN] cannot get group member id from input: %v", id)
	}
	return project, tag, err
}
