package gitlab

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Fourcast/go-gitlab"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceGitlabProjectShareGroup() *schema.Resource {
	acceptedAccessLevels := []string{"guest", "reporter", "developer", "maintainer"}

	return &schema.Resource{
		Create: resourceGitlabProjectShareGroupCreate,
		Read:   resourceGitlabProjectShareGroupRead,
		Delete: resourceGitlabProjectShareGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
			},
			"access_level": {
				Type:         schema.TypeString,
				ValidateFunc: validateValueFunc(acceptedAccessLevels),
				ForceNew:     true,
				Required:     true,
			},
		},
	}
}

func resourceGitlabProjectShareGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	groupId := d.Get("group_id").(int)
	projectId := d.Get("project_id").(string)
	accessLevelId := accessLevelID[d.Get("access_level").(string)]

	options := &gitlab.ShareWithGroupOptions{
		GroupID:     &groupId,
		GroupAccess: &accessLevelId,
	}
	log.Printf("[DEBUG] create gitlab project membership for %d in %s", options.GroupID, projectId)

	_, err := client.Projects.ShareProjectWithGroup(projectId, options)
	if err != nil {
		return err
	}
	groupIdString := strconv.Itoa(groupId)
	d.SetId(buildTwoPartID(&projectId, &groupIdString))
	return resourceGitlabProjectShareGroupRead(d, meta)
}

func resourceGitlabProjectShareGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	id := d.Id()
	log.Printf("[DEBUG] read gitlab project projectMember %s", id)

	projectId, groupId, err := projectIdAndGroupIdFromId(id)
	if err != nil {
		return err
	}

	projectInformation, _, err := client.Projects.GetProject(projectId, nil)
	if err != nil {
		log.Printf("[DEBUG] failed to read gitlab project %s: %s", id, err)
		d.SetId("")
		return nil
	}

	for _, v := range projectInformation.SharedWithGroups {
		if groupId == v.GroupID {
			resourceGitlabProjectShareGroupSetToState(d, v, &projectId)
		}
	}

	return nil
}

func projectIdAndGroupIdFromId(id string) (string, int, error) {
	projectId, groupIdString, err := parseTwoPartID(id)
	if err != nil {
		return "", 0, fmt.Errorf("Error parsing ID: %s", id)
	}

	groupId, err := strconv.Atoi(groupIdString)
	if err != nil {
		return "", 0, fmt.Errorf("Can not determine group id: %v", id)
	}

	return projectId, groupId, nil
}

func resourceGitlabProjectShareGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	id := d.Id()
	projectId, groupId, err := projectIdAndGroupIdFromId(id)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Delete gitlab project membership %v for %s", groupId, projectId)

	_, err = client.Projects.DeleteSharedProjectFromGroup(projectId, groupId)
	return err
}

func resourceGitlabProjectShareGroupSetToState(d *schema.ResourceData, group struct {
	GroupID          int    "json:\"group_id\""
	GroupName        string "json:\"group_name\""
	GroupAccessLevel int    "json:\"group_access_level\""
}, projectId *string) {

	//This cast is needed due to an inconsistency in the upstream API
	//GroupAcessLevel is returned as an int but the map we lookup is sorted by the int alias AccessLevelValue
	convertedAccessLevel := gitlab.AccessLevelValue(group.GroupAccessLevel)

	d.Set("project_id", projectId)
	d.Set("group_id", group.GroupID)
	d.Set("access_level", accessLevel[convertedAccessLevel])

	groupId := strconv.Itoa(group.GroupID)
	d.SetId(buildTwoPartID(projectId, &groupId))
}
