package provider

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

const encoding = "base64"

// NOTE: this lock is a bit of a hack to prevent parallel calls to the GitLab Repository Files API.
//
//	If it is called concurrently, the API will return a 400 error along the lines of:
//	```
//	(400 Bad Request) DELETE https://gitlab.com/api/v4/projects/30716/repository/files/somefile.yaml: 400
//	{message: 9:Could not update refs/heads/master. Please refresh and try again..}
//	```
//
//	This lock only solves half of the problem, where the provider is responsible for
//	the concurrency. The other half is if the API is called outside of terraform at the same time
//	this resource makes calls to the API.
//	To mitigate this, simple retries are used.
var resourceGitlabRepositoryFileApiLock = newLock()

var _ = registerResource("gitlab_repository_file", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_repository_file`" + ` resource allows to manage the lifecycle of a file within a repository.

-> **Timeouts** Default timeout for *Create*, *Update* and *Delete* is one minute and can be configured in the ` + "`timeouts`" + ` block.

-> **Implementation Detail** GitLab is unable to handle concurrent calls to the GitLab repository files API for the same project.
   Therefore, this resource queues every call to the repository files API no matter of the project, which may slow down the terraform
   execution time for some configurations. In addition, retries are performed in case a refresh is required because another application
   changed the repository at the same time.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/repository_files.html)`,

		CreateContext: resourceGitlabRepositoryFileCreate,
		ReadContext:   resourceGitlabRepositoryFileRead,
		UpdateContext: resourceGitlabRepositoryFileUpdate,
		DeleteContext: resourceGitlabRepositoryFileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		// the schema matches https://docs.gitlab.com/ee/api/repository_files.html#create-new-file-in-repository
		// However, we don't support the `encoding` parameter as it seems to be broken.
		// Only a value of `base64` is supported, all others, including the documented default `text`, lead to
		// a `400 {error: encoding does not have a valid value}` error.
		Schema: constructSchema(
			map[string]*schema.Schema{
				"branch": {
					Description: "Name of the branch to which to commit to.",
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
				},
				"commit_message": {
					Description: "Commit message.",
					Type:        schema.TypeString,
					Required:    true,
				},
				"start_branch": {
					Description: "Name of the branch to start the new commit from.",
					Type:        schema.TypeString,
					Optional:    true,
				},
				"author_email": {
					Description: "Email of the commit author.",
					Type:        schema.TypeString,
					Optional:    true,
				},
				"author_name": {
					Description: "Name of the commit author.",
					Type:        schema.TypeString,
					Optional:    true,
				},
			},
			gitlabRepositoryFileGetSchema(),
		),
	}
})

func resourceGitlabRepositoryFileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	project := d.Get("project").(string)
	filePath := d.Get("file_path").(string)

	log.Printf("[DEBUG] gitlab_repository_file: waiting for lock to create %s/%s", project, filePath)
	if err := resourceGitlabRepositoryFileApiLock.lock(ctx); err != nil {
		return diag.FromErr(err)
	}
	defer resourceGitlabRepositoryFileApiLock.unlock()
	log.Printf("[DEBUG] gitlab_repository_file: got lock to create %s/%s", project, filePath)

	client := meta.(*gitlab.Client)
	// NOTE: for backwards-compatibility reasons, we also support an already given base64 encoding,
	//       otherwise we encode the `content` to base64.
	content := d.Get("content").(string)
	if _, err := base64.StdEncoding.DecodeString(content); err != nil {
		log.Printf("[DEBUG] gitlab_repository_file: given content '%s' is not a valid base64 encoded string, encoding it ...", content)
		content = base64.StdEncoding.EncodeToString([]byte(content))
	}

	options := &gitlab.CreateFileOptions{
		Branch:        gitlab.String(d.Get("branch").(string)),
		Encoding:      gitlab.String(encoding),
		AuthorEmail:   gitlab.String(d.Get("author_email").(string)),
		AuthorName:    gitlab.String(d.Get("author_name").(string)),
		Content:       gitlab.String(content),
		CommitMessage: gitlab.String(d.Get("commit_message").(string)),
	}
	if startBranch, ok := d.GetOk("start_branch"); ok {
		options.StartBranch = gitlab.String(startBranch.(string))
	}
	if executeFilemode, ok := d.GetOk("execute_filemode"); ok {
		options.ExecuteFilemode = gitlab.Bool(executeFilemode.(bool))
	}

	err := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		repositoryFile, _, err := client.RepositoryFiles.CreateFile(project, filePath, options, gitlab.WithContext(ctx))
		if err != nil {
			if isRefreshError(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		d.SetId(resourceGitLabRepositoryFileBuildId(project, repositoryFile.Branch, repositoryFile.FilePath))
		return nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGitlabRepositoryFileRead(ctx, d, meta)
}

func resourceGitlabRepositoryFileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project, branch, filePath, err := resourceGitLabRepositoryFileParseId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	options := &gitlab.GetFileOptions{
		Ref: gitlab.String(branch),
	}

	repositoryFile, _, err := client.RepositoryFiles.GetFile(project, filePath, options, gitlab.WithContext(ctx))
	if err != nil {
		if strings.Contains(err.Error(), "404 File Not Found") {
			log.Printf("[WARN] file %s not found, removing from state", filePath)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	configContent := d.Get("content").(string)
	log.Printf("[DEBUG] gitlab_repository_file: comparing content of %s with %s", repositoryFile.Content, configContent)
	// NOTE: for backwards-compatibility reasons, we also support an already given base64 encoding,
	//       otherwise we encode the `content` to base64.
	if _, err := base64.StdEncoding.DecodeString(configContent); err != nil {
		// if `content` is config is not a base64 encoded string, we decode the one from the API, too
		// in case it's base64 encoded, else we don't decode it.
		if decodedContent, err := base64.StdEncoding.DecodeString(repositoryFile.Content); err == nil {
			repositoryFile.Content = string(decodedContent)
		}
	}

	d.SetId(resourceGitLabRepositoryFileBuildId(project, branch, repositoryFile.FilePath))
	d.Set("branch", repositoryFile.Ref)
	stateMap := gitlabRepositoryFileToStateMap(project, repositoryFile)
	if err = setStateMapInResourceData(stateMap, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGitlabRepositoryFileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	project, branch, filePath, err := resourceGitLabRepositoryFileParseId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] gitlab_repository_file: waiting for lock to update %s/%s", project, filePath)
	if err := resourceGitlabRepositoryFileApiLock.lock(ctx); err != nil {
		return diag.FromErr(err)
	}
	defer resourceGitlabRepositoryFileApiLock.unlock()
	log.Printf("[DEBUG] gitlab_repository_file: got lock to update %s/%s", project, filePath)

	client := meta.(*gitlab.Client)

	readOptions := &gitlab.GetFileOptions{
		Ref: gitlab.String(branch),
	}

	// NOTE: for backwards-compatibility reasons, we also support an already given base64 encoding,
	//       otherwise we encode the `content` to base64.
	content := d.Get("content").(string)
	if _, err := base64.StdEncoding.DecodeString(content); err != nil {
		log.Printf("[DEBUG] gitlab_repository_file: given content '%s' is not a valid base64 encoded string, encoding it ...", content)
		content = base64.StdEncoding.EncodeToString([]byte(content))
	}

	updateOptions := &gitlab.UpdateFileOptions{
		Branch:        gitlab.String(branch),
		Encoding:      gitlab.String(encoding),
		AuthorEmail:   gitlab.String(d.Get("author_email").(string)),
		AuthorName:    gitlab.String(d.Get("author_name").(string)),
		Content:       gitlab.String(content),
		CommitMessage: gitlab.String(d.Get("commit_message").(string)),
	}
	if startBranch, ok := d.GetOk("start_branch"); ok {
		updateOptions.StartBranch = gitlab.String(startBranch.(string))
	}
	if executeFilemode, ok := d.GetOk("execute_filemode"); ok {
		updateOptions.ExecuteFilemode = gitlab.Bool(executeFilemode.(bool))
	}

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		// NOTE: we also re-read the file to obtain an eventually changed `LastCommitID` for which we needed the refresh
		existingRepositoryFile, _, err := client.RepositoryFiles.GetFile(project, filePath, readOptions, gitlab.WithContext(ctx))
		if err != nil {
			return resource.NonRetryableError(err)
		}

		updateOptions.LastCommitID = gitlab.String(existingRepositoryFile.LastCommitID)
		_, _, err = client.RepositoryFiles.UpdateFile(project, filePath, updateOptions, gitlab.WithContext(ctx))
		if err != nil {
			if isRefreshError(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGitlabRepositoryFileRead(ctx, d, meta)
}

func resourceGitlabRepositoryFileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	project, branch, filePath, err := resourceGitLabRepositoryFileParseId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] gitlab_repository_file: waiting for lock to delete %s/%s", project, filePath)
	if err := resourceGitlabRepositoryFileApiLock.lock(ctx); err != nil {
		return diag.FromErr(err)
	}
	defer resourceGitlabRepositoryFileApiLock.unlock()
	log.Printf("[DEBUG] gitlab_repository_file: got lock to delete %s/%s", project, filePath)

	client := meta.(*gitlab.Client)

	readOptions := &gitlab.GetFileOptions{
		Ref: gitlab.String(branch),
	}
	deleteOptions := &gitlab.DeleteFileOptions{
		Branch:        gitlab.String(d.Get("branch").(string)),
		AuthorEmail:   gitlab.String(d.Get("author_email").(string)),
		AuthorName:    gitlab.String(d.Get("author_name").(string)),
		CommitMessage: gitlab.String(fmt.Sprintf("[DELETE]: %s", d.Get("commit_message").(string))),
	}

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		// NOTE: we also re-read the file to obtain an eventually changed `LastCommitID` for which we needed the refresh

		existingRepositoryFile, _, err := client.RepositoryFiles.GetFile(project, filePath, readOptions, gitlab.WithContext(ctx))
		if err != nil {
			return resource.NonRetryableError(err)
		}

		deleteOptions.LastCommitID = gitlab.String(existingRepositoryFile.LastCommitID)
		resp, err := client.RepositoryFiles.DeleteFile(project, filePath, deleteOptions, gitlab.WithContext(ctx))
		if err != nil {
			if isRefreshError(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(fmt.Errorf("%s failed to delete repository file: (%s) %v", d.Id(), resp.Status, err))
		}
		return nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGitLabRepositoryFileParseId(id string) (string, string, string, error) {
	parts := strings.SplitN(id, ":", 3)
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("Unexpected ID format (%q). Expected project:branch:repository_file_path", id)
	}

	return parts[0], parts[1], parts[2], nil
}

func resourceGitLabRepositoryFileBuildId(project string, branch string, filePath string) string {
	return fmt.Sprintf("%s:%s:%s", project, branch, filePath)
}

func isRefreshError(err error) bool {
	var httpErr *gitlab.ErrorResponse
	return errors.As(err, &httpErr) &&
		httpErr.Response.StatusCode == http.StatusBadRequest &&
		strings.Contains(httpErr.Message, "Please refresh and try again")
}
