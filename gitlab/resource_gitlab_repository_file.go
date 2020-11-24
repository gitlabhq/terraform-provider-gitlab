package gitlab

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func resourceGitlabRepositoryFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabRepositoryFileCreate,
		Read:   resourceGitlabRepositoryFileRead,
		Update: resourceGitlabRepositoryFileUpdate,
		Delete: resourceGitlabRepositoryFileDelete,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"file": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"branch": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"author_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"author_email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"commit_message": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceGitlabRepositoryFileCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	file := d.Get("file").(string)

	options := &gitlab.CreateFileOptions{
		Branch:        gitlab.String(d.Get("branch").(string)),
		AuthorEmail:   gitlab.String(d.Get("author_email").(string)),
		AuthorName:    gitlab.String(d.Get("author_name").(string)),
		Content:       gitlab.String(d.Get("content").(string)),
		CommitMessage: gitlab.String(d.Get("commit_message").(string)),
		Encoding:      gitlab.String("base64"),
	}

	repositoryFile, _, err := client.RepositoryFiles.CreateFile(project, file, options)
	if err != nil {
		return err
	}

	d.SetId(repositoryFile.FilePath)

	return resourceGitlabRepositoryFileRead(d, meta)
}

func resourceGitlabRepositoryFileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	fileID := d.Id()
	options := &gitlab.GetFileOptions{
		Ref: gitlab.String(d.Get("branch").(string)),
	}

	repositoryFile, _, err := client.RepositoryFiles.GetFile(project, fileID, options)
	if err != nil {
		return err
	}

	d.Set("project", project)
	d.Set("file", repositoryFile.FileName)
	d.Set("content", repositoryFile.Content)
	d.Set("branch", repositoryFile.Ref)

	return nil
}

func resourceGitlabRepositoryFileUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	file := d.Get("file").(string)
	options := &gitlab.UpdateFileOptions{
		Branch:        gitlab.String(d.Get("branch").(string)),
		AuthorEmail:   gitlab.String(d.Get("author_email").(string)),
		AuthorName:    gitlab.String(d.Get("author_name").(string)),
		Content:       gitlab.String(d.Get("content").(string)),
		CommitMessage: gitlab.String(d.Get("commit_message").(string)),
		Encoding:      gitlab.String("base64"),
		//TODO: add LastCommitID
	}

	if d.HasChange("branch") {
		options.Branch = gitlab.String(d.Get("branch").(string))
	}

	if d.HasChange("author_email") {
		options.AuthorEmail = gitlab.String(d.Get("author_email").(string))
	}

	if d.HasChange("author_name") {
		options.AuthorName = gitlab.String(d.Get("author_name").(string))
	}

	if d.HasChange("content") {
		options.Content = gitlab.String(d.Get("content").(string))
	}

	if d.HasChange("commit_message") {
		options.CommitMessage = gitlab.String(d.Get("commit_message").(string))
	}
	_, _, err := client.RepositoryFiles.UpdateFile(project, file, options)
	if err != nil {
		return err
	}

	return resourceGitlabRepositoryFileRead(d, meta)
}

func resourceGitlabRepositoryFileDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	file := d.Get("file").(string)
	options := &gitlab.DeleteFileOptions{
		Branch:        gitlab.String(d.Get("branch").(string)),
		AuthorEmail:   gitlab.String(d.Get("author_email").(string)),
		AuthorName:    gitlab.String(d.Get("author_name").(string)),
		CommitMessage: gitlab.String(d.Get("commit_message").(string)),
		//TODO: add LastCommitID
	}

	resp, err := client.RepositoryFiles.DeleteFile(project, file, options)
	if err != nil {
		return fmt.Errorf("%s failed to delete repository file: (%s) %v", d.Id(), resp.Status, err)
	}
	return err
}
