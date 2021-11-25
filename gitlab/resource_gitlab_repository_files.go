package gitlab

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/xanzy/go-gitlab"
	"net/http"
	"strings"
)

func resourceGitLabRepositoryFiles() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabRepositoryFilesCreate,
		Read:   resourceGitlabRepositoryFilesRead,
		Update: resourceGitlabRepositoryFilesUpdate,
		Delete: resourceGitlabRepositoryFilesDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				s := strings.Split(d.Id(), ":")
				if len(s) != 3 {
					d.SetId("")
					return nil, fmt.Errorf("invalid Repository File import format; expected '{project_id}:{branch}:{file_path1,file_path2,file_path3}'")
				}
				project, branch, filePathsStr := s[0], s[1], s[2]

				var filePathsBuilder strings.Builder
				var flattenedFiles []map[string]string
				for _, filePath := range strings.Split(filePathsStr, ",") {
					filePathsBuilder.WriteString(filePath)
					flattenedFiles = append(flattenedFiles, map[string]string{
						"file_path": filePath,
					})
				}

				d.Set("file", flattenedFiles)
				d.SetId(fmt.Sprintf("%d", schema.HashString(filePathsBuilder.String())))
				d.Set("project", project)
				d.Set("branch", branch)

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"branch": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"start_branch": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"author_email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"author_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"commit_message": {
				Type:     schema.TypeString,
				Required: true,
			},
			"file": {
				Type:     schema.TypeSet,
				Required: true,
				Set:      fileHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"content": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceGitlabRepositoryFilesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	var filePathsBuilder strings.Builder
	var actions []*gitlab.CommitActionOptions
	files := d.Get("file").(*schema.Set).List()
	for _, file := range files {
		fileMap := file.(map[string]interface{})
		actions = append(actions, &gitlab.CommitActionOptions{
			Action:   gitlab.FileAction(gitlab.FileCreate),
			FilePath: gitlab.String(fileMap["file_path"].(string)),
			Content:  gitlab.String(fileMap["content"].(string)),
		})
		filePathsBuilder.WriteString(fileMap["file_path"].(string))
	}
	opts := &gitlab.CreateCommitOptions{
		Branch:        gitlab.String(d.Get("branch").(string)),
		AuthorEmail:   gitlab.String(d.Get("author_email").(string)),
		AuthorName:    gitlab.String(d.Get("author_name").(string)),
		CommitMessage: gitlab.String(d.Get("commit_message").(string)),
		Actions:       actions,
	}
	if startBranch, ok := d.GetOk("start_branch"); ok {
		opts.StartBranch = gitlab.String(startBranch.(string))
	}

	_, _, err := client.Commits.CreateCommit(d.Get("project").(string), opts)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", schema.HashString(filePathsBuilder.String())))
	return resourceGitlabRepositoryFilesRead(d, meta)
}

func resourceGitlabRepositoryFilesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	files := d.Get("file").(*schema.Set).List()

	var flattenedFiles []map[string]string
	for _, fileRaw := range files {
		fileMap := fileRaw.(map[string]interface{})
		file, resp, err := client.RepositoryFiles.GetFile(d.Get("project").(string), fileMap["file_path"].(string), &gitlab.GetFileOptions{
			Ref: gitlab.String(d.Get("branch").(string)),
		})
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				continue
			}
			return err
		}
		content, err := base64.StdEncoding.DecodeString(file.Content)
		if err != nil {
			return err
		}
		flattenedFiles = append(flattenedFiles, map[string]string{
			"file_path": file.FilePath,
			"content":   string(content),
		})
	}

	d.Set("file", flattenedFiles)
	d.Set("project", d.Get("project").(string))

	return nil
}

func resourceGitlabRepositoryFilesUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	if d.HasChange("file") {
		desiredStateRaw, actualStateRaw := d.GetChange("file")
		desiredState, actualState := desiredStateRaw.(*schema.Set), actualStateRaw.(*schema.Set)

		filesToCreate := actualState.Difference(desiredState)
		filesToDelete := desiredState.Difference(actualState)

		// If a file_path is in filesToCreate and filesToDelete, then it should be removed from both of them and placed in filesToUpdate.
		// This happens when the file content changes on remote.
		var filesToUpdate []interface{}
		for _, fileToCreateRaw := range filesToCreate.List() {
			fileToCreate := fileToCreateRaw.(map[string]interface{})
			if fileToCreateFilePath, ok := fileToCreate["file_path"]; ok {

				for _, fileToDeleteRaw := range filesToDelete.List() {
					fileToDelete := fileToDeleteRaw.(map[string]interface{})
					if fileToDeleteFilePath, ok := fileToDelete["file_path"]; ok {

						if fileToCreateFilePath.(string) == fileToDeleteFilePath.(string) {
							filesToUpdate = append(filesToUpdate, map[string]interface{}{
								"file_path": fileToCreateFilePath,
								"content":   fileToCreate["content"].(string),
							})

							filesToCreate.Remove(fileToCreateRaw)
							filesToDelete.Remove(fileToDeleteRaw)
						}
					}
				}
			}
		}

		var actions []*gitlab.CommitActionOptions
		createActions := func(action *gitlab.FileActionValue, rawFiles []interface{}) {
			for _, rawFile := range rawFiles {
				rawFileMap := rawFile.(map[string]interface{})
				actions = append(actions, &gitlab.CommitActionOptions{
					Action:   action,
					FilePath: gitlab.String(rawFileMap["file_path"].(string)),
					Content:  gitlab.String(rawFileMap["content"].(string)),
				})
			}
		}
		createActions(gitlab.FileAction(gitlab.FileCreate), filesToCreate.List())
		createActions(gitlab.FileAction(gitlab.FileDelete), filesToDelete.List())
		createActions(gitlab.FileAction(gitlab.FileUpdate), filesToUpdate)

		opts := &gitlab.CreateCommitOptions{
			Branch:        gitlab.String(d.Get("branch").(string)),
			AuthorEmail:   gitlab.String(d.Get("author_email").(string)),
			AuthorName:    gitlab.String(d.Get("author_name").(string)),
			CommitMessage: gitlab.String(d.Get("commit_message").(string)),
			Actions:       actions,
		}
		if startBranch, ok := d.GetOk("start_branch"); ok {
			opts.StartBranch = gitlab.String(startBranch.(string))
		}

		_, _, err := client.Commits.CreateCommit(d.Get("project").(string), opts)
		if err != nil {
			return err
		}
	}

	return resourceGitlabRepositoryFilesRead(d, meta)
}

func resourceGitlabRepositoryFilesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	files := d.Get("file").(*schema.Set).List()

	var actions []*gitlab.CommitActionOptions
	for _, file := range files {
		actions = append(actions, &gitlab.CommitActionOptions{
			Action:   gitlab.FileAction(gitlab.FileDelete),
			FilePath: gitlab.String(file.(map[string]interface{})["file_path"].(string)),
		})
	}
	opts := &gitlab.CreateCommitOptions{
		Branch:        gitlab.String(d.Get("branch").(string)),
		AuthorEmail:   gitlab.String(d.Get("author_email").(string)),
		AuthorName:    gitlab.String(d.Get("author_name").(string)),
		CommitMessage: gitlab.String(d.Get("commit_message").(string)),
		Actions:       actions,
	}
	if startBranch, ok := d.GetOk("start_branch"); ok {
		opts.StartBranch = gitlab.String(startBranch.(string))
	}

	_, _, err := client.Commits.CreateCommit(d.Get("project").(string), opts)
	if err != nil {
		return err
	}

	return resourceGitlabRepositoryFilesRead(d, meta)
}

func fileHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	buf.WriteString(fmt.Sprintf("%s,", m["file_path"].(string)))
	buf.WriteString(fmt.Sprintf("%s,", m["content"].(string)))

	return schema.HashString(buf.String())
}
