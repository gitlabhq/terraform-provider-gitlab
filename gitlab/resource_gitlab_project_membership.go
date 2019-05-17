package gitlab

import (
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func resourceGitlabProjectMembership() *schema.Resource {
	acceptedAccessLevels := make([]string, 0, len(accessLevelID))
	for k := range accessLevelID {
		if k != "owner" {
			acceptedAccessLevels = append(acceptedAccessLevels, k)
		}
	}
	return &schema.Resource{
		Create: resourceGitlabProjectMembershipCreate,
		Read:   resourceGitlabProjectMembershipRead,
		Update: resourceGitlabProjectMembershipUpdate,
		Delete: resourceGitlabProjectMembershipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"user_id": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
			},
			"access_level": {
				Type:         schema.TypeString,
				ValidateFunc: validateValueFunc(acceptedAccessLevels),
				Required:     true,
			},
		},
	}
}

func resourceGitlabProjectMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	userId := d.Get("user_id").(int)
	projectId := d.Get("project_id").(string)
	accessLevelId := accessLevelID[d.Get("access_level").(string)]

	options := &gitlab.AddProjectMemberOptions{
		UserID:      &userId,
		AccessLevel: &accessLevelId,
	}
	log.Printf("[DEBUG] create gitlab project membership for %d in %s", options.UserID, projectId)

	_, _, err := client.ProjectMembers.AddProjectMember(projectId, options)
	if err != nil {
		return err
	}
	userIdString := strconv.Itoa(userId)
	d.SetId(buildTwoPartID(&projectId, &userIdString))
	return resourceGitlabProjectMembershipRead(d, meta)
}

func resourceGitlabProjectMembershipRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	id := d.Id()
	log.Printf("[DEBUG] read gitlab project projectMember %s", id)

	projectId, userId, e := projectIdAndUserIdFromId(id)
	if e != nil {
		return e
	}

	projectMember, _, err := client.ProjectMembers.GetProjectMember(projectId, userId)
	if err != nil {
		return err
	}

	resourceGitlabProjectMembershipSetToState(d, projectMember, &projectId)
	return nil
}

func projectIdAndUserIdFromId(id string) (string, int, error) {
	projectId, userIdString, err := parseTwoPartID(id)
	userId, e := strconv.Atoi(userIdString)
	if err != nil {
		e = err
	}
	if e != nil {
		log.Printf("[WARN] cannot get project member id from input: %v", id)
	}
	return projectId, userId, e
}

func resourceGitlabProjectMembershipUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	userId := d.Get("user_id").(int)
	projectId := d.Get("project_id").(string)
	accessLevelId := accessLevelID[strings.ToLower(d.Get("access_level").(string))]

	options := gitlab.EditProjectMemberOptions{
		AccessLevel: &accessLevelId,
	}
	log.Printf("[DEBUG] update gitlab project membership %v for %s", userId, projectId)

	_, _, err := client.ProjectMembers.EditProjectMember(projectId, userId, &options)
	if err != nil {
		return err
	}
	return resourceGitlabProjectMembershipRead(d, meta)
}

func resourceGitlabProjectMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	id := d.Id()
	projectId, userId, e := projectIdAndUserIdFromId(id)
	if e != nil {
		return e
	}

	log.Printf("[DEBUG] Delete gitlab project membership %v for %s", userId, projectId)

	_, err := client.ProjectMembers.DeleteProjectMember(projectId, userId)
	return err
}

func resourceGitlabProjectMembershipSetToState(d *schema.ResourceData, projectMember *gitlab.ProjectMember, projectId *string) {

	d.Set("project_id", projectId)
	d.Set("user_id", projectMember.ID)
	d.Set("access_level", accessLevel[projectMember.AccessLevel])

	userId := strconv.Itoa(projectMember.ID)
	d.SetId(buildTwoPartID(projectId, &userId))
}
