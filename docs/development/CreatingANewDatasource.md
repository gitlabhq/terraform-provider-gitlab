---
page_title: "Creating a New Datasource Using the Terraform Plugin Framework"
---

# Creating a New Datasource Using the Terraform Plugin Framework

This tutorial is to help new contributors out when creating a new datasource.
It will walk through a step-by-step guide of creating a new datasource using the [Terraform Provider Framework](https://developer.hashicorp.com/terraform/plugin/framework).
This guide will assume that a development environment has already been set up by following the `Developing The Provider` section of the CONTRIBUTING.md documentation.

<!-- Use "yzhang.markdown-all-in-one" plugin to keep this up to date in vscode -->
- [Creating a New Datasource Using the Terraform Plugin Framework](#creating-a-new-datasource-using-the-terraform-plugin-framework)
  - [Step 1: Understand the API from GitLab](#step-1-understand-the-api-from-gitlab)
  - [Step 2: Create the Datasource struct](#step-2-create-the-datasource-struct)
  - [Step 3: Create the `Metadata` function](#step-3-create-the-metadata-function)
  - [Step 4: Create the Schema](#step-4-create-the-schema)
  - [Step 5: Create the `Configure` function](#step-5-create-the-configure-function)
  - [Step 6: Create the `Read` function](#step-6-create-the-read-function)
  - [Step 7: Verify the datasource is properly structured](#step-7-verify-the-datasource-is-properly-structured)
  - [Step 8: Create documentation](#step-8-create-documentation)
  - [Step 9: Create tests](#step-9-create-tests)
  - [Conclusion](#conclusion)

## Step 1: Understand the API from GitLab

When creating a new datasource, the GitLab terraform provider follows the [Terraform Provider Best Practices](https://developer.hashicorp.com/terraform/plugin/best-practices/hashicorp-provider-design-principles) whenever possible.
This means that a new datasource meets a couple of criteria:

1. One datasource aligns as closely to one Read API as possible.
2. The attributes of the datasource align to the attributes of the underlying API.

The [`datasource_gitlab_cluster_agent`](../internal/provider/datasource_gitlab_cluster_agent.go) datasource will be used as a step-by-step example.
This datasource aligns to the [Get Details About An Agent API](https://docs.gitlab.com/api/cluster_agents/#get-details-about-an-agent) exposed by GitLab.
When creating a datasource, first ensure that the relevant APIs are present in GitLab.
If it's not clear whether an API exists for a datasource, create an issue on the GitLab Terraform Provider project and ask.

## Step 2: Create the Datasource struct

In the Terraform Plugin framework, each datasource is represented by a struct that implements one or more interfaces.
For the sake of keeping this tutorial simple, these interfaces won't be covered in details.
Each datasource is created within its own `go` file, named `datasource_<datasource_name>.go`; in this case, `datasource_gitlab_cluster_agent.go`.

```golang
type gitlabClusterAgentDataSource struct {
    client *gitlab.Client // This is required for making calls to GitLab later
}
```

## Step 3: Create the `Metadata` function

This function tells terraform the name of the datasource:

```golang
func (d *gitlabClusterAgentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_cluster_agent"
}
```

## Step 4: Create the Schema

The schema for the datasource handles multiple responsibilities:

1. It ensures that the input data is the correct type (`number` vs `string`).
2. It ensures that the input data is properly validated (matches any validation rules).
3. It ensures that the input data has all the necessarily required fields.

Copy all the required and optional attributes from the GitLab API into a schema struct.
This struct represents the attributes that a user can use to configure the datasource:

```golang
type gitlabClusterAgentDataSourceModel struct {
    ID              types.String `tfsdk:"id"`
    Project         types.String `tfsdk:"project"`
    Name            types.String `tfsdk:"name"`
    AgentID         types.Int64  `tfsdk:"agent_id"`
    CreatedAt       types.String `tfsdk:"created_at"`
    CreatedByUserID types.Int64  `tfsdk:"created_by_user_id"`
}
```

There are a couple of things to notice about this struct:

1. The types for each attribute use the `types` package. This is because `types.String` can have a nil value, whereas a primitive `string` cannot.
2. The `tfsdk` tag value maps to the string value in the schema.
3. Each new struct like this must have a unique name. The terraform provider uses the naming convention of `gitlab<datasourceName><resource type, either Resource or DataSource>Model`. That means a cluster agent data source would be named `gitlabClusterAgentDataSourceModel`.

With the schema struct and the datasource struct created, it's time to start implementing the datasource functions.

The first function to create is the `Schema` function, which defines a `schema.Schema` struct representing the schema and all the validations required for the datasource.

```golang
func (d *gitlabClusterAgentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: `The ` + "`gitlab_cluster_agent`" + ` data source retrieves details about a GitLab Agent for Kubernetes.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/api/cluster_agents/)`,
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                MarkdownDescription: "The ID of this data source. In the format <project:agent_id>",
                Computed:            true,
            },
            "project": schema.StringAttribute{
                MarkdownDescription: "ID or full path of the project maintained by the authenticated user.",
                Required:            true,
            },
            "name": schema.StringAttribute{
                MarkdownDescription: "The Name of the agent.",
                Computed:            true,
            },
            "agent_id": schema.Int64Attribute{
                MarkdownDescription: "The ID of the agent.",
                Required:            true,
            },
            "created_at": schema.StringAttribute{
                MarkdownDescription: "The ISO8601 datetime when the agent was created.",
                Computed:            true,
            },
            "created_by_user_id": schema.Int64Attribute{
                MarkdownDescription: "The ID of the user who created the agent.",
                Computed:            true,
            },
        },
    }
}
```

Similar to the schema struct above, there are a couple things to take note of in the above `Schema` func.

1. The `Schema` func itself is part of the `datasource.DataSource` interface. Make sure it has the proper inputs!
2. Each `Schema` must have a `MarkdownDescription`. This will appear in the terraform documentation on the provider's site. It must include a link to the relevant REST API documentation.
3. Each `Schema` must have a `Attributes` map, which contains a minimum of one `schema.Attribute` in its map. This map is where plan-time validation happens. Within each `schema.Attribute`, several key properties are required:
  a `MarkdownDescription` if the documentation that will appear on the terraform documentation site for that attribute.
  b `Required` denotes whether the attribute is required for the datasource. Datasources missing required attributes will fail at plan-time.
  c `Computed` denotes whether the datasource will compute values for that attribute that may differ from the plan. If `Computed` is set to `true`, then storing a value that's different from the terraform config won't result in a diff being identified unless the value is explicitly set in the config.
  d `Validators` accepts validator functions that can be used to validate inputs at plan time.

For more information on various properties of the schema attributes, read the [Terraform Plugin Framework Schema Documentation](https://developer.hashicorp.com/terraform/plugin/framework/handling-data/schemas).

## Step 5: Create the `Configure` function

After the schema function has been written, the `Configure` function needs to be written.

```golang
// Configure adds the provider configured client to the datasource.
func (d *gitlabClusterAgentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
    // Prevent panic if the provider has not been configured.
    if req.ProviderData == nil {
        return
    }

    r.client = req.ProviderData.(*gitlab.Client)
}
```

This function will be nearly identical on every datasource.
The logic sets the client in the datasource struct to be the value configured in the provider.
This ensures that when making calls from the `r.Client` that they're authenticated and configured properly.

## Step 6: Create the `Read` function

The `Read` function is called any time terraform `refresh` is called, either by a `plan`, an `apply`, or an explicit `refresh`.

```golang
func (d *gitlabClusterAgentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data gitlabClusterAgentDataSourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    project := data.Project.ValueString()
    agentID := int(data.AgentID.ValueInt64())

    agent, _, err := d.client.ClusterAgents.GetAgent(project, agentID, gitlab.WithContext(ctx))
    if err != nil {
        resp.Diagnostics.AddError("Failed to get cluster agent", err.Error())
        return
    }

    agentIDStr := strconv.Itoa(agentID)
    data.ID = types.StringValue(utils.BuildTwoPartID(&project, &agentIDStr))
    data.Project = types.StringValue(project)
    data.Name = types.StringValue(agent.Name)
    data.AgentID = types.Int64Value(int64(agent.ID))
    data.CreatedAt = types.StringValue(agent.CreatedAt.Format(time.RFC3339))
    data.CreatedByUserID = types.Int64Value(int64(agent.CreatedByUserID))
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
```

There are several things to notice in the `Read` function:

1. `resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)` will read all the attributes from the request, and store them in `data`. This
allows the data object to be used in downstream calls to retrieve data from the config in a typesafe manner.
2. `if resp.Diagnostics.HasError() {return}` checks to ensure that reading the config didn't encounter an error, and exits before any changes
are made if an error was encountered.
3. `d.client.ClusterAgents.GetAgent(..)` calls GitLab's API to get the cluster agent data.
4. The cluster agent data is then set into the schema struct.
5. Finally, once the data object has been updated, `resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)` is called to set the data back into the state file.

This may look complicated, but it's following three steps:

1. Read the attributes from the configuration, which should contain a unique identifier for the datasource
2. Use the unique identifier to query the datasource from the API
3. Store the response from the API back into the state

All `Read` functions will follow this same pattern.

## Step 7: Verify the datasource is properly structured

All the functions created in the above tutorial are created because they are required by interfaces.
If the datasource has been structured appropriately, and has the appropriate functions, then the datasource can be assigned to the interfaces successfully.
To verify this, the following assignments should be added to the datasource implementation:

```golang
var (
    _ datasource.DataSource              = &gitlabClusterAgentDataSource{} //requires `Schema` and `Read` function
    _ datasource.DataSourceWithConfigure = &gitlabClusterAgentDataSource{} //requires the `Configure` function
)
```

If any functions are missing, a compile error will provide details around which functions are missing.
There will currently be one compilation error with `Datasource`, noting that an `init()` function is required to register the datasource with the provider.
Adding the `init` function will require only several lines of code to complete the datasource:

```golang
func init() {
    registerDataSource(NewGitlabClusterAgentDataSource)
}

func NewGitlabClusterAgentDataSource() datasource.DataSource {
    return &gitlabClusterAgentDataSource{}
}
```

This will return a new instance of the struct created earlier in the tutorial, and register it to the provider which will allow its use.

## Step 8: Create documentation

Now that the datasource has been created and the code has been written, we need to document our datasource.
Most of the documentation will be generated by the `make generate` command, which will read the documentation from the `Schema` block created in step 4.
However, one additional piece of documentation is required - an example of the datasource in terraform hcl.

To provide these examples, create a new folder under `/examples/datasources`.
The folder should have the same name as the terraform datasource.
In this example, that means `/examples/datasources/gitlab_cluster_agent`.

To provide the datasource example, create a `data-source.tf` file, containing an example configuration, and any comments that may help an end user.

```hcl
data "gitlab_cluster_agent" "example" {
    project  = "12345"
    agent_id = 1
}

```

Finally, run `make generate` from the root of the repository.
This make take about 20-30 seconds to install all go dependencies the first time it runs, but it will generate all documents within the `/docs` folder that are necessary for the datasource.

## Step 9: Create tests

Every datasource should have tests associated with it, and test principles can be located in the [CONTRIBUTING.md file](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/blob/main/CONTRIBUTING.md).
Tests are created in a separate `go` file, using a standard naming convention.
Append `_test` to the end of your datasource's file name.
For example, `gitlab_cluster_agent` is in `datasource_gitlab_cluster_agent.go`, so the tests are located in `datasource_gitlab_cluster_agent_test.go`.
To ensure that test logic is kept separate from the provider, build tags are used for acceptance tests:

```golang
//go:build acceptance
// +build acceptance
```

These tags will ensure that logic contained in the test files doesn't get compiled into the provider's binary.
They also ensure that the tests are run when the acceptance test CI jobs are invoked.

Creating a test for a datasource involves using the terraform testing framework, and creating a `resource.TestCase`.
This test will run terraform commands in the order specified by the test steps, and will execute check methods after each step.
At the end of each test case, `terraform destroy` will be run, then a function specified in `CheckDestroy`.
Here is an example:

```golang
func TestAccDataSourceGitlabClusterAgent_basic(t *testing.T) {
    testProject := testutil.CreateProject(t)
    testAgent := testutil.CreateClusterAgents(t, testProject.ID, 1)[0]

    resource.ParallelTest(t, resource.TestCase{
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: fmt.Sprintf(`
                    data "gitlab_cluster_agent" "this" {
                        project           = "%d"
                        agent_id          = %d
                    }
                    `, testProject.ID, testAgent.ID,
                ),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr("data.gitlab_cluster_agent.this", "name", testAgent.Name),
                    resource.TestCheckResourceAttr("data.gitlab_cluster_agent.this", "created_at", testAgent.CreatedAt.Format(time.RFC3339)),
                    resource.TestCheckResourceAttr("data.gitlab_cluster_agent.this", "created_by_user_id", fmt.Sprintf("%d", testAgent.CreatedByUserID)),
                ),
            },
        },
    })
}
```

Key items to notice about the above example include:

1. the `Steps` attribute of the `TestCase` include a slice of `TestStep`. The first step (and any non-import steps) includes a `Config`
attribute that specifies what the terraform configuration is. This will run a `terraform apply` to create that datasource.
2. The `Check` attribute includes a set of `CheckFunc`, and the `resource` package provides a set of implementations that can be used
for checking things like values, or to check that an attribute is set without checking its value.

If any TestStep check functions fail, the Test will fail and produce an error.

These steps are usually covered in a `_basic` test.
In a complicated datasource, many tests may be required to fully cover a new datasource.
It's always better to err on the side of creating too many tests than it is to create not enough.

## Conclusion

During this tutorial, a new `gitlab_cluster_agent` datasource has been created, including the schema, Read operation, and more.
Hopefully this tutorial has been helpful, but every tutorial has room for improvement.
If there are improvements to any examples, or if any steps are confusing, please feel free to open an MR to continue to iterate on the documentation.

Thank you for taking the time to read through this tutorial, and the whole community looks forward to working with you on making the gitlab terraform provider excellent!
