package gitlab

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

const encoding = "base64"

func resourceGitLabRepositoryFile() *schema.Resource {
	return &schema.Resource{
		Description: "This resource allows you to create and manage GitLab repository files.\n\n" +
			"**Limitations**:\n\n" +
			"The [GitLab Repository Files API](https://docs.gitlab.com/ee/api/repository_files.html)\n" +
			"can only create, update or delete a single file at the time.\n" +
			"The API will also\n" +
			"[fail with a `400`](https://docs.gitlab.com/ee/api/repository_files.html#update-existing-file-in-repository)\n" +
			"response status code if the underlying repository is changed while the API tries to make changes.\n" +
			"Therefore, it's recommended to make sure that you execute it with\n" +
			"[`-parallelism=1`](https://www.terraform.io/docs/cli/commands/apply.html#parallelism-n)\n" +
			"and that no other entity than the terraform at hand makes changes to the\n" +
			"underlying repository while it's executing.",

		CreateContext: resourceGitlabRepositoryFileCreate,
		ReadContext:   resourceGitlabRepositoryFileRead,
		UpdateContext: resourceGitlabRepositoryFileUpdate,
		DeleteContext: resourceGitlabRepositoryFileDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		// the schema matches https://docs.gitlab.com/ee/api/repository_files.html#create-new-file-in-repository
		// However, we don't support the `encoding` parameter as it seems to be broken.
		// Only a value of `base64` is supported, all others, including the documented default `text`, lead to
		// a `400 {error: encoding does not have a valid value}` error.
		Schema: map[string]*schema.Schema{
			"project": {
				Description: "The ID of the project.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"file_path": {
				Description: "The full path of the file. It must be relative to the root of the project without a leading slash `/`.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"branch": {
				Description: "Name of the branch to which to commit to.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
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
			"content": {
				Description:  "base64 encoded file content. No other encoding is currently supported, because of a [GitLab API bug](https://gitlab.com/gitlab-org/gitlab/-/issues/342430).",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateBase64Content,
			},
			"commit_message": {
				Description: "Commit message.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"encoding": {
				Description: "Content encoding.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceGitlabRepositoryFileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	filePath := d.Get("file_path").(string)

	options := &gitlab.CreateFileOptions{
		Branch:        gitlab.String(d.Get("branch").(string)),
		Encoding:      gitlab.String(encoding),
		AuthorEmail:   gitlab.String(d.Get("author_email").(string)),
		AuthorName:    gitlab.String(d.Get("author_name").(string)),
		Content:       gitlab.String(d.Get("content").(string)),
		CommitMessage: gitlab.String(d.Get("commit_message").(string)),
	}
	if startBranch, ok := d.GetOk("start_branch"); ok {
		options.StartBranch = gitlab.String(startBranch.(string))
	}

	repositoryFile, _, err := client.RepositoryFiles.CreateFile(project, filePath, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resourceGitLabRepositoryFileBuildId(project, repositoryFile.Branch, repositoryFile.FilePath))
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

	d.SetId(resourceGitLabRepositoryFileBuildId(project, branch, repositoryFile.FilePath))
	d.Set("project", project)
	d.Set("file_path", repositoryFile.FilePath)
	d.Set("branch", repositoryFile.Ref)
	d.Set("encoding", repositoryFile.Encoding)
	d.Set("content", repositoryFile.Content)

	return nil
}

func resourceGitlabRepositoryFileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project, branch, filePath, err := resourceGitLabRepositoryFileParseId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	readOptions := &gitlab.GetFileOptions{
		Ref: gitlab.String(branch),
	}

	existingRepositoryFile, _, err := client.RepositoryFiles.GetFile(project, filePath, readOptions, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	options := &gitlab.UpdateFileOptions{
		Branch:        gitlab.String(branch),
		Encoding:      gitlab.String(encoding),
		AuthorEmail:   gitlab.String(d.Get("author_email").(string)),
		AuthorName:    gitlab.String(d.Get("author_name").(string)),
		Content:       gitlab.String(d.Get("content").(string)),
		CommitMessage: gitlab.String(d.Get("commit_message").(string)),
		LastCommitID:  gitlab.String(existingRepositoryFile.LastCommitID),
	}
	if startBranch, ok := d.GetOk("start_branch"); ok {
		options.StartBranch = gitlab.String(startBranch.(string))
	}

	_, _, err = client.RepositoryFiles.UpdateFile(project, filePath, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGitlabRepositoryFileRead(ctx, d, meta)
}

func resourceGitlabRepositoryFileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project, branch, filePath, err := resourceGitLabRepositoryFileParseId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	readOptions := &gitlab.GetFileOptions{
		Ref: gitlab.String(branch),
	}

	existingRepositoryFile, _, err := client.RepositoryFiles.GetFile(project, filePath, readOptions, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	options := &gitlab.DeleteFileOptions{
		Branch:        gitlab.String(d.Get("branch").(string)),
		AuthorEmail:   gitlab.String(d.Get("author_email").(string)),
		AuthorName:    gitlab.String(d.Get("author_name").(string)),
		CommitMessage: gitlab.String(fmt.Sprintf("[DELETE]: %s", d.Get("commit_message").(string))),
		LastCommitID:  gitlab.String(existingRepositoryFile.LastCommitID),
	}

	resp, err := client.RepositoryFiles.DeleteFile(project, filePath, options)
	if err != nil {
		return diag.Errorf("%s failed to delete repository file: (%s) %v", d.Id(), resp.Status, err)
	}

	return nil
}

func validateBase64Content(v interface{}, k string) (we []string, errors []error) {
	content := v.(string)
	if _, err := base64.StdEncoding.DecodeString(content); err != nil {
		errors = append(errors, fmt.Errorf("given repository file content '%s' is not base64 encoded, but must be", content))
	}
	return
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
