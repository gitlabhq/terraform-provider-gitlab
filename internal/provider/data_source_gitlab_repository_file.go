package provider

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_repository_file", func() *schema.Resource {
	return &schema.Resource{
		Description: "Allows you to receive information about file in repository like name, size, content. File content is Base64 encoded. This endpoint can be accessed without authentication if the repository is publicly accessible.",
		ReadContext: dataSourceGitlabRepositoryFileRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Description: "The ID of the project.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"file_path": {
				Description: "The full path of the file. It must be relative to the root of the project without a leading slash `/`.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"ref": {
				Description: "The name of branch, tag or commit.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"file_name": {
				Description: "String, file name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"size": {
				Description: "Integer, file size.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"encoding": {
				Description: "String, file encoding.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"content": {
				Description: "String, base64 encoded file content.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"content_sha256": {
				Description: "String, content sha256 digest.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"blob_id": {
				Description: "String, blob id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"commit_id": {
				Description: "String, commit id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_commit_id": {
				Description: "String, last commit id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
})

func dataSourceGitlabRepositoryFileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	filePath := d.Get("file_path").(string)

	options := &gitlab.GetFileOptions{
		Ref: gitlab.String(d.Get("ref").(string)),
	}

	repositoryFile, resp, err := client.RepositoryFiles.GetFile(project, filePath, options, gitlab.WithContext(ctx))
	if err != nil {
		log.Printf("[DEBUG] file %s not found, response %v", filePath, resp)
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s:%s:%s", project, repositoryFile.Ref, repositoryFile.FilePath))

	d.Set("project", project)
	d.Set("file_name", repositoryFile.FileName)
	d.Set("file_path", repositoryFile.FilePath)
	d.Set("size", repositoryFile.Size)
	d.Set("encoding", repositoryFile.Encoding)
	d.Set("content", repositoryFile.Content)
	d.Set("content_sha256", repositoryFile.SHA256)
	d.Set("ref", repositoryFile.Ref)
	d.Set("blob_id", repositoryFile.BlobID)
	d.Set("commit_id", repositoryFile.CommitID)
	d.Set("last_commit_id", repositoryFile.LastCommitID)

	return nil
}
