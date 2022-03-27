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
		Description: `The ` + "`gitlab_repository_file`" + ` data source allows details of a file in a repository to be retrieved.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/repository_files.html)`,

		ReadContext: dataSourceGitlabRepositoryFileRead,
		Schema:      datasourceSchemaFromResourceSchema(gitlabRepositoryFileGetSchema(), []string{"project", "file_path", "ref"}, nil),
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

	stateMap := gitlabRepositoryFileToStateMap(project, repositoryFile)
	if err := setStateMapInResourceData(stateMap, d); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
