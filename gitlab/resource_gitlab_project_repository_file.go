package gitlab

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabProjectRepositoryFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabProjectRepositoryFileCreate,
		Read:   resourceGitlabProjectRepositoryFileRead,
		Update: resourceGitlabProjectRepositoryFileUpdate,
		Delete: resourceGitlabProjectRepositoryFileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"path": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"branch": {
				Type:     schema.TypeString,
				Required: true,
				// ForceNew: true,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"commit_message": {
				Type:     schema.TypeString,
				Required: true,
			},
			"delete_commit_message": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sha256": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGitlabProjectRepositoryFileCreate(d *schema.ResourceData, meta interface{}) error {
	project := d.Get("project").(string)
	path := d.Get("path").(string)
	branch := d.Get("branch").(string)
	content := d.Get("content").(string)
	commitMessage := d.Get("commit_message").(string)

	options := gitlab.CreateFileOptions{
		Branch:        &branch,
		Content:       &content,
		CommitMessage: &commitMessage,
	}

	log.Printf("[DEBUG] Project %s create gitlab repository file %+v", project, options)

	client := meta.(*gitlab.Client)

	fileInfo, _, err := client.RepositoryFiles.CreateFile(project, path, &options)
	if err != nil {
		return err
	}

	d.SetId(buildThreePartID(&project, &fileInfo.Branch, &fileInfo.FilePath))

	return resourceGitlabProjectRepositoryFileRead(d, meta)
}

func resourceGitlabProjectRepositoryFileRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] read gitlab repository file %s", d.Id())

	projectID, branch, path, err := parseThreePartID(d.Id())
	if err != nil {
		return err
	}
	d.Set("project", projectID)

	client := meta.(*gitlab.Client)

	file, resp, err := client.RepositoryFiles.GetFile(projectID, path, &gitlab.GetFileOptions{
		Ref: &branch,
	})
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error getting repository file %q on project %q, branch %q: %v", path, projectID, branch, err)
	}

	d.Set("path", file.FilePath)
	d.Set("branch", file.Ref)
	content, err := base64.StdEncoding.DecodeString(file.Content)
	if err != nil {
		return fmt.Errorf("error decoding content of file %q on project %q, branch %q: %v", path, projectID, branch, err)
	}
	d.Set("content", content)
	d.Set("sha256", file.SHA256)

	return nil
}

func resourceGitlabProjectRepositoryFileUpdate(d *schema.ResourceData, meta interface{}) error {
	projectID, branch, path, err := parseThreePartID(d.Id())
	if err != nil {
		return err
	}

	content := d.Get("content").(string)
	commitMessage := d.Get("commit_message").(string)

	options := gitlab.UpdateFileOptions{
		Branch:        &branch,
		Content:       &content,
		CommitMessage: &commitMessage,
	}

	log.Printf("[DEBUG] Project %s update gitlab repository file %q on branch %q", projectID, path, branch)

	client := meta.(*gitlab.Client)

	if _, _, err := client.RepositoryFiles.UpdateFile(projectID, path, &options); err != nil {
		return fmt.Errorf("error updating gitlab repository file %q on project %q, branch %q: %v", path, projectID, branch, err)
	}

	return resourceGitlabProjectRepositoryFileRead(d, meta)
}

func resourceGitlabProjectRepositoryFileDelete(d *schema.ResourceData, meta interface{}) error {
	project, branch, path, err := parseThreePartID(d.Id())
	if err != nil {
		return err
	}

	commitMessage := d.Get("commit_message").(string)
	deleteCommitMessage, ok := d.GetOk("delete_commit_message")
	if ok {
		commitMessage = deleteCommitMessage.(string)
	} else {
		commitMessage = fmt.Sprintf("deleted: %s", commitMessage)
	}

	options := gitlab.DeleteFileOptions{
		Branch:        &branch,
		CommitMessage: &commitMessage,
	}

	log.Printf("[DEBUG] Project %s delete gitlab repository file %q on branch %q", project, path, branch)

	client := meta.(*gitlab.Client)

	_, err = client.RepositoryFiles.DeleteFile(project, path, &options)
	if err != nil {
		return err
	}

	return nil
}
