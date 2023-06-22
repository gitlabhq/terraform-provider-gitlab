---
page_title: "Creating a New Resource Using the Terraform Plugin Framework"
---

# Creating a New Resource Using the Terraform Plugin Framework

This tutorial is meant to help new contributors out when creating new resource. It will walk through a 
step-by-step guide of creating a new resource using the 
[Terraform Provider Framework](https://developer.hashicorp.com/terraform/plugin/framework),
since that is how all new resources are added to the GitLab terraform provider, as noted in the
[CONTRIBUTING.md](/CONTRIBUTING.md). This guide will assume that a development environment has already
been set up by following the `Developing The Provider` section of the CONTRIBUTING.md documentation.

<!-- Use "yzhang.markdown-all-in-one" plugin to keep this up to date in vscode -->
- [Creating a New Resource Using the Terraform Plugin Framework](#creating-a-new-resource-using-the-terraform-plugin-framework)
	- [Step 1: Understand the API from GitLab](#step-1-understand-the-api-from-gitlab)
	- [Step 2: Create the Resource struct](#step-2-create-the-resource-struct)
	- [Step 3: Create the Schema](#step-3-create-the-schema)
	- [Step 4: Create the `Config` function](#step-4-create-the-config-function)
	- [Step 5: Create the CRUD operations](#step-5-create-the-crud-operations)
		- [Step 5a: Create the `Read` function](#step-5a-create-the-read-function)
		- [Step 5b: Create the `Create` function](#step-5b-create-the-create-function)
	- [Step 6: Add import support](#step-6-add-import-support)
	- [Step 7: Verify the resource is properly structured](#step-7-verify-the-resource-is-properly-structured)
	- [Step 8: Create documentation](#step-8-create-documentation)
	- [Step 8: Create tests](#step-8-create-tests)
- [Conclusion](#conclusion)


## Step 1: Understand the API from GitLab

When creating a new resource, the GitLab terraform provider follows the
[Terraform Provider Best Practices](https://developer.hashicorp.com/terraform/plugin/best-practices/hashicorp-provider-design-principles)
whenever possible. This means that a new resource meets a couple of criteria:

1. One resource aligns as closely to one set of CRUD APIs as possible.
2. The attributes of the resource align to the attributes of the underlying APIs.

For this example, the [`resource_gitlab_application`](../internal/provider/resource_gitlab_application.go)
resource will be used as a step-by-step example. This resource aligns to the 
[Applications API](https://docs.gitlab.com/ee/api/applications.html) exposed by GitLab. When creating
a resource, first ensure that the relevant APIs are present in GitLab. If it's not clear whether an
api exists for a resource, create an issue on the GitLab Terraform Provider project and ask!

## Step 2: Create the Resource struct

In the Terraform Plugin framework, each resource is represented by a struct that implements one or more
interfaces. For the sake of keeping this tutorial simple, these interfaces won't be covered in details. However,
creating the resource struct will be the first step in creating a new resource. Each resource is created
within its own `go` file, named `resource_<resource_name>.go`; in this case, `resource_gitlab_application.go`. 

```golang
type gitlabApplicationResource struct {
	client *gitlab.Client // This is required for making calls to GitLab later
}
```

## Step 3: Create the Schema

The schema for the resource handles multiple responsibilities during `terraform plan` and `terraform apply`:

1. It ensures that the input data is the correct type (`number` vs `string`).
2. It ensures that the input data is properly validated (matches any validation rules).
3. It ensures that the input data has all the necessarily required fields.

As a result, the schema is the natural starting point for creating a resource. The best place to start
for creating a resource is to copy all the required and optional attributes from the GitLab API into the
schema struct. To define the schema for the resource, first, create a struct representing the attributes
that a user can use to configure the resource:

```golang
type gitlabApplicationResourceModel struct {
	Name         types.String `tfsdk:"name"`
	RedirectURL  types.String `tfsdk:"redirect_url"`
	Scopes       types.Set    `tfsdk:"scopes"`
	Confidential types.Bool   `tfsdk:"confidential"`

	Id            types.String `tfsdk:"id"`
	Secret        types.String `tfsdk:"secret"`
	ApplicationId types.String `tfsdk:"application_id"`
}
```
 
There are a couple of things to notice about this struct:

1. The types for each attribute use the `types` package. This is because `types.String` can have a nil value,
whereas a primative `string` cannot.
2. The `tfsdk` tag value maps to the string value in the schema.
3. Each new struct like this must have a unique name. The terraform provider uses the naming convention of
`gitlab<resourceName><resource type, either Resource or Data>Model`. That means an application data source 
would be named `gitlabApplicationDataModel`.

After the schema struct is created, the next step is to create a second struct representing the resource itself. This
struct will then implement all the functions that are required for performing terraform CRUD (Create, Read,
Update, Delete) operations.

```golang
type gitlabApplicationResource struct {
	client *gitlab.Client
}
```

This struct is very simple, and just accepts a client reference. This client will be used to make REST calls to
the GitLab instance configured in the provider.

With the schema struct and the resource struct created, it's time to start implementing the resource functions.

The first function to create is the `Schema` function, which defines a `schema.Schema` struct representing the schema
and all the validations required for the resource. The schema block is very large, so the full block will not be copied here. 
The full schema function can be read 
[in the repository, linked here](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/blob/main/internal/provider/resource_gitlab_application.go#L63)

```golang
func (r *gitlabApplicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf(`The ` + "`gitlab_application`" + ` resource allows to manage the lifecycle of applications in gitlab.

~> In order to use a user for a user to create an application, they must have admin priviledges at the instance level.
To create an OIDC application, a scope of "openid".

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/applications.html)`),

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of this Terraform resource. In the format of `<application_id>`.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the application.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			// additional schema resources past this point.
		}
	}
}
```

Similar to the schema struct above, there are a couple things to take note of in the above `Schema` func.

1. The `Schema` func itself is part of the `resource.Resource` interface. Make sure it has the proper inputs!
2. Each `Schema` must have a `MarkdownDescription`. This will appear in the terraform documentation on the provider's site.
3. Each `Schema` must have a `Attributes` map, which contains a minimum of one `schema.Attribute` in its map. This map
is where plan-time validation happens. Within each `schema.Attribute`, several key properties are required:
  - `MarkdownDescription` if the documentation that will appear on the terraform documentation site for that attribute.
  - `Required` denotes whether the attribute is required for the resource. Resources missing required attribute will fail at plan-time.
  - `Computed` denotes whether the resource will compute values for that attribute that may differ from the plan. If `Computed` is 
  set to `true`, then storing a value that's different from the terraform config won't result in a diff being identified unless the 
  value is explicitly set in the config. 
  - `Validators` accepts validator functions that can be used to validate inputs at plan time.
  - `PlanModifiers` accepts modifier functions that can change how the resource identifies plan changes.

For more information on various properties of the schema attributes, feel free to read the 
[Terraform Plugin Framework Schema Documentation](https://developer.hashicorp.com/terraform/plugin/framework/handling-data/schemas).

## Step 4: Create the `Config` function

After the schema function has been written, the `Config` function needs to be written. Don't worry, this one is much easier!

```golang
// Configure adds the provider configured client to the resource.
func (r *gitlabApplicationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*gitlab.Client)
}
```

This function will be nearly identical on every resource. The logic simply sets the client in the resource struct to be the value 
configured in the provider. This ensures that when making calls from the `r.Client` that they're authenticated and configured properly.

## Step 5: Create the CRUD operations

### Step 5a: Create the `Read` function

Finally, it's time to create the CRUD functions for the resource. The CRUD functions (Create, Read, Update, and Delete) are responsible
for using the `r.Client` to make the changes to the GitLab instance. Terraform will automatically call the correct function based on 
the terraform plan that's generated before the apply:

- If a resource is labelled as `create`, the `Create` function will be called. 
- If a resource is labelled as `update`, the `Update` function will be called.
- If a resource is labelled as `destroy`, the `Delete` function will be called. 
- The `Read` function is called any time terraform `refresh` is called, either by a `plan`, an `apply`, or an explicit `refresh`.

Creating a CRUD funcion involves reading the attributes from the terraform configuration, then passing them to the API call necessary
to manipulate the resource in GitLab. This document will demonstrate a `Read` and `Create` function; other functions can be read from
the `gitlab_application_settings.go` file.

First, creating the `Read` function:

```golang
func (r *gitlabApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *gitlabApplicationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	application, err := findGitlabApplication(r.client, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("GitLab API error occurred", fmt.Sprintf("Unable to create application: %s", err.Error()))
		return
	}

	tflog.Trace(ctx, "found application", map[string]interface{}{
		"application": gitlab.Stringify(application),
	})

	r.applicationModelToState(application, data)
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func findGitlabApplication(client *gitlab.Client, desiredId string) (*gitlab.Application, error) {

	options := gitlab.ListApplicationsOptions{
		PerPage: 20,
		Page:    1,
	}

	for options.Page != 0 {
		paginatedApplications, resp, err := client.Applications.ListApplications(&options)
		if err != nil {
			return nil, fmt.Errorf("unable to list applications. %s", err)
		}

		for i := range paginatedApplications {
			if strconv.Itoa(paginatedApplications[i].ID) == desiredId {
				return paginatedApplications[i], nil
			}
		}

		options.Page = resp.NextPage
	}

	// if we loop through the pages and haven't found it, we should error
	return nil, fmt.Errorf("unable to find application with id: %s", desiredId)
}

func (r *gitlabApplicationResource) applicationModelToState(application *gitlab.Application, data *gitlabApplicationResourceModel) {
	// need to check this
	// For reads, the secret will be empty, in which case we shouldn't set the state
	if application.Secret != "" {
		data.Secret = types.StringValue(application.Secret)
	}
	data.Id = types.StringValue(strconv.Itoa(application.ID))
	data.Confidential = types.BoolValue(application.Confidential)
	data.Name = types.StringValue(application.ApplicationName)
	data.RedirectURL = types.StringValue(application.CallbackURL)
	data.ApplicationId = types.StringValue(application.ApplicationID)
}

```
Like before, there are several things to notice in the `Read` function:

1. `resp.Diagnostics.Append(req.State.Get(ctx, &data)...)` will read all the attributes from the request, and store them in `data`. This
allows the data object to be used in downstream calls to retrieve data from the config in a typesafe manner.
2. `if resp.Diagnostics.HasError() {return}` checks to ensure that reading the config didn't encounter an error, and exits before any changes
are made if an error was encountered.
3. `findGitlabApplication` demonstrates that CRUD functions can invoke helper functions in Go, just like any other Go code. `findGitlabApplication`
is used to paginate through listed applications and find the application ID from the config.
4. If an application was returned from the API, `applicationModelToState` is invoked to set all the properties of the `data` object. Notice the
`types.StringValue` calls being made, which take a string primitive and convert it to a `types.String` object.
5. Finally, once the data object has been updated, `resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)` is called to set the data back into 
the state file.

This may look complicated, but it's following three steps:

1. Read the attributes from the configuration, which should contain a unique identifier for the resource
2. Use the unique identifier to query the resource from the API
3. Store the response from the API back into the state

All `Read` functions will follow this same pattern.

### Step 5b: Create the `Create` function

Creating a resource, similarly, involves reading from the configuration, populating an API call, and then storting the created resource
back into state.

```golang
// Create creates a new upstream resources and adds it into the Terraform state.
func (r *gitlabApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *gitlabApplicationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "creating application", map[string]interface{}{
		"scopes": data.Scopes.String(),
	})
	scopes := conv.StringSetToStrings(data.Scopes)
	if resp.Diagnostics.HasError() {
		return
	}

	formatted_scopes := strings.Join(scopes, " ")

	// configure GitLab API call
	options := &gitlab.CreateApplicationOptions{
		Name:        gitlab.String(data.Name.ValueString()),
		RedirectURI: gitlab.String(data.RedirectURL.ValueString()),
		Scopes:      gitlab.String(formatted_scopes),
	}

	if !data.Confidential.IsNull() {
		options.Confidential = gitlab.Bool(data.Confidential.ValueBool())
	}

	// Create application
	application, _, err := r.client.Applications.CreateApplication(options)
	if err != nil {
		resp.Diagnostics.AddError("GitLab API error occurred", fmt.Sprintf("Unable to create application: %s", err.Error()))
		return
	}

	r.applicationModelToState(application, data)
	// Log the creation of the resource
	tflog.Debug(ctx, "created an application", map[string]interface{}{
		"name": data.Name.ValueString(), "id": data.Id.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
```

Things to notice about this function:

1. Similar to the `Read` function, the `Create` function starts with `var data *gitlabApplicationResourceModel` to read the config
into the `data` struct. This is used through the rest of the function to reference the config.
2. the API accepts a list of `scopes` which is separated by a space. Since the terraform config is going to accept a formatted list
object, we need to format that []string into a single string separated by spaces. That's what `formatted_scopes := strings.Join(scopes, " ")`
is doing. This is an example where the terraform provider may format inputs slightly before passing them to the API to ensure that
the input follows terraform best practices.
3. `options := &gitlab.CreateApplicationOptions{...}` sets all the required values into an options struct. Any values not configured
will not be passed to the API, and either defaults will be used or the value will not be set by the API.
4. `if !data.Confidential.IsNull() {...}` checks to see if an optional value is null. If the value is not null, it's added to the `options`
struct to be passed to the API.
5. `application, _, err := r.client.Applications.CreateApplication(options)` calls the API with the input options to create the application.
6. Just like in the `Read` func, `r.applicationModelToState(application, data)` sets the newly created application into the data object
7. Finally, like in the `Read` func, `resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)` sets the data into state.

## Step 6: Add import support

Most terraform resources should support the ability to use `import` to load pre-existing resources into the terraform state. To do this,
a function called `ImportState` is added to the resource:

```golang
func (r *gitlabApplicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
```

Most resources will use the `ImportStatePassthroughID` function to handle import. The critical piece of logic in the above example is 
the uses of `path.Root("id")`, which specifies that the `id` attribute is the primary key for the resource. When a user calls
`terraform import <resource_path> <resource_id>`, the value of `resource_id` will be set into the attribute specified in `path.Root()`.

Then, the first time `terraform refresh` is called (either by explicitly calling it, or via a `plan` or `apply` operation), terraform
will execute the `Read` function, and read the resource specified by that `id`.

## Step 7: Verify the resource is properly structured

All the functions created in the above tutorial are created because they are required by interfaces. If the resource has been structured
appropriately, and has the appropriate functions (remember, an `Update` and `Delete` function are required in addition to the ones
created in this toturial) then the resource can be assigned to the interfaces successfully. To verify this, the following assignements
should be added to the resource implementation:

```golang
var (
	_ resource.Resource                = &gitlabApplicationResource{} //requires `Schema` and CRUD functions
	_ resource.ResourceWithConfigure   = &gitlabApplicationResource{} //requires the `Config` function
	_ resource.ResourceWithImportState = &gitlabApplicationResource{} //requires the `ImportState` function
)
```

If any functions are missing, a compile error will provide details around which functions are missing. There will currently be one
compilation error with `Resource`, noting that an `init()` function is required to register the resource with the provider.  Adding
the `init` function is easy, and will require only several lines of code to complete the resource:

```golang
func init() {
	registerResource(NewGitLabApplicationResource)
}

// NewGitLabApplicationResource is a helper function to simplify the provider implementation.
func NewGitLabApplicationResource() resource.Resource {
	return &gitlabApplicationResource{}
}
```

This will return a new instance of the struct created earlier in the tutorial, and register it to the provider which will allow its 
use.


## Step 8: Create documentation

Now that the resource has been created and the code has been written, we need to document our resources. Most of the documentation
will be generated by the `make generate` command, which will read the documentation from the `Schema` block created in step 3. However,
two additional pieces of documentation are required:

1. An example of the resource in terraform hcl
2. An example of how to import the resource

To provide these examples, create a new folder under `/examples/resources`; the folder should have the same name as the terraform resource.
In this tutorial, that means `gitlab_application`.

To provide the `import` example, create an `import.sh` file, containing any logic needed to import the resource, and any comments that
may help an end user.

```bash
# Gitlab applications can be imported with their id, e.g.
terraform import gitlab_application.example "1"
```

To provide the resource example, create a `resource.tf` file, containing an example configuration, and any comments that may help an end user.

```hcl
resource "gitlab_application" "oidc" {
  confidential = true
  scopes       = ["openid"]
  name         = "company_oidc"
  redirect_url = "https://mycompany.com"
}
```

Finally, run `make generate` from the root of the repository. This make take about 20-30 seconds to install all go dependencies the first time
it runs, but it will generate all documents within the `/docs` folder that are necessary for the resource.

## Step 8: Create tests

Every resource should have tests associated with it, and test principles can be located in the 
[CONTRIBUTING.md file](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/blob/main/CONTRIBUTING.md). Tests are created in a separate
`go` file, using a standard naming convention, appending `_test` to the end of your resource's file name. For example, the `gitlab_application`
is in `resource_gitlab_application.go`, so the tests are located in `resource_gitlab_application_test.go`. To ensure that test logic is 
kept separate from the provider, build tags are used for acceptance tests:

```golang
//go:build acceptance
// +build acceptance
```
These tags will ensure that logic contained in the test files doesn't get compiled into the provider's binary, and they ensure that the tests
are run properly when the acceptance test CI jobs are invoked.

Creating a test for a resource involves using the terraform testing framework, and creating a `resource.TestCase`. This test will run terraform
commands in the order specified by the test steps, and will execute check methods after each step. At the end of each test case, `terraform destroy`
will be run, then a function specified in `CheckDestroy`. Here is an example:

```golang
func TestAcc_GitlabApplication_basic(t *testing.T) {
	name := acctest.RandString(10)
	url := "https://my_website.com"

	resource.ParallelTest(t, resource.TestCase{

		ProtoV6ProviderFactories: testAccProtoV6MuxProviderFactories,
		CheckDestroy:             testAcc_GitlabApplication_CheckDestroy(),

		Steps: []resource.TestStep{
			// Create a basic application.
			{
				Config: fmt.Sprintf(`
				resource "gitlab_application" "this" {
					name     = %q
					redirect_url = %q
					scopes = ["openid"]
					confidential = true
				}`, name, url),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gitlab_application.this", "redirect_url", url),
					resource.TestCheckResourceAttr("gitlab_application.this", "scopes.0", "openid"),
				),
			},
			// Verify upstream attributes with an import.
			{
				ResourceName:            "gitlab_application.this",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret", "scopes"},
			},
		},
	})
}

func testAcc_GitlabApplication_CheckDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {

			if rs.Type == "gitlab_application" {
				application, err := findGitlabApplication(testutil.TestGitlabClient, rs.Primary.ID)
				if err == nil {
					return fmt.Errorf("Found GitLab application that should have been deleted: %s", gitlab.Stringify(application))
				}
			}
		}
		return nil
	}
}
```

Key items to notice about the above example include:
1. the `Steps` attribute of the `TestCase` include a slice of `TestStep`. The first step (and any non-import steps) includes a `Config`
attribute that specifies what the terraform configuration is. This will run a `terraform apply` to create that resource.
2. The `Check` attribute includes a set of `CheckFunc`, and the `resource` package provides a set of implementations that can be used
for checking things like values, or to check that an attribute is set without checking its value.
3. In the second `TestStep`, the `ImportState` attribute is set to true. This will run `terraform import` and import the resource
specified in the `ResourceName` attribute. If any attributes don't match the value returned from the import commany, this `TestStep` 
will return an error. Since the `secret` and `scopes` values are not returned from the API, those cannot be imported, so those two
attributes are ignored by including them in the `ImportStateVerifyIgnore` attribute.
4. The `CheckDestroy` attribute accepts a function. This function loops over the state values until it identified the `gitlab_application`
resource. It retrieves the `id` of the resource (which is the primary key of the application) then attempts to retrieve that application
using the GitLab API. If the application is still present, it means the application didn't delete properly, and the function returns an
`error` object.

If any TestStep check functions fail, any imports fail, or any destroy functions fail, the Test will fail and produce an error.

All new resources are expected to test a minimum of 4 operations:
1. Create a new resource
2. Update an existing resource (usually the one created in step 1)
3. Import the resource
4. Destroy the resource

These steps are usually covered in a `_basic` test. In a complicated resource, many tests may be required to fully cover a new resource.
It's always better to err on the side of creating too many tests than it is to create not enough.

# Conclusion

During this tutorial, a new `gitlab_application` resource has been created, including the schema, CRUD operations, Import operations, and more.
Hopefully this tutorial has been helpful, but every tutorial has room for improvement. If there are improvements to any examples, or if any steps
are confusing, please feel free to open an MR are continue to iterate on the documentation.

Thank you for taking the time to read through this tutorial, and the whole community looks forward to working with you on making the gitlab terraform 
provider excellent!