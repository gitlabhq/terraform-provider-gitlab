<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Terraform Provider for Gitlab
=============================

- [Documentation](https://www.terraform.io/docs/providers/gitlab/index.html)
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)
- Build status:
  - [![Unit Tests](https://github.com/gitlabhq/terraform-provider-gitlab/workflows/Unit%20Tests/badge.svg?branch=master)](https://github.com/gitlabhq/terraform-provider-gitlab/actions?query=workflow%3A%22Unit+Tests%22+branch%3Amaster)
  - [![Acceptance Tests](https://github.com/gitlabhq/terraform-provider-gitlab/workflows/Acceptance%20Tests/badge.svg?branch=master)](https://github.com/gitlabhq/terraform-provider-gitlab/actions?query=workflow%3A%22Acceptance+Tests%22+branch%3Amaster)
  - ![Website Build](https://github.com/gitlabhq/terraform-provider-gitlab/workflows/Website%20Build/badge.svg?branch=master)

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) >= 1.14 (to build the provider plugin)

## Developing The Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.14+ is *required*).

1. Clone the git repository.

   ```sh
   $ git clone git@github.com:gitlabhq/terraform-provider-gitlab
   $ cd terraform-provider-gitlab
   ```

2. Build the provider with `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

   ```sh
   $ make build
   ```

### Running Tests

The acceptance tests can run against a Gitlab instance where you have a token with administrator permissions (likely not gitlab.com).

#### Option 1: Run tests against a local Gitlab container with docker-compose

This option is the easiest and requires [docker-compose](https://docs.docker.com/compose/install/) (version 1.13+) to be installed on your machine.

1. Start the Gitlab container. It will take about 5 minutes for the container to become healthy.

  ```sh
  $ make testacc-up
  ```

2. Run the acceptance tests. The full suite takes 10-20 minutes to run.

  ```sh
  $ make testacc
  ```

3. Stop the Gitlab container.

  ```sh
  $ make testacc-down
  ```

#### Option 2: Run tests against your own Gitlab instance

If you have your own hosted Gitlab instance, you can run the tests against it directly.

```sh
$ make testacc GITLAB_TOKEN=example123 GITLAB_BASE_URL=https://example.com/api/v4
```

`GITLAB_TOKEN` must be a valid token for an account with admin privileges.

#### Testing Tips

* **Gitlab Community Edition and Gitlab Enterprise Edition:**

  This module supports both Gitlab CE and Gitlab EE. We run tests on Gitlab EE,
  but can't run them on pull requests from forks.

  Features that only work on one flavour can use the following helpers as
  SkipFunc: `isRunningInEE` and `isRunningInCE`. You can see an example of this
  for [gitlab_project_level_mr_approvals](gitlab/resource_gitlab_project_level_mr_approvals_test.go)
  tests.

* **Run EE tests:**

  If you have a `Gitlab-license.txt` you can run Gitlab EE, which will enable the full suite of tests:

  ```sh
  $ make testacc-up SERVICE=gitlab-ee
  ```

* **Run a single test:**

  You can pass a pattern to the `RUN` variable to run a reduced number of tests. For example:

  ```sh
  $ make testacc RUN=TestAccGitlabGroup
  ```

   ...will run all tests for the `gitlab_group` resource.

* **Debug a test in an IDE:**

  First start the Gitlab container with `make testacc-up`.
  Then run the desired Go test as you would normally from your IDE, but configure your run configuration to set these environment variables:

  ```
  GITLAB_TOKEN=ACCTEST
  GITLAB_BASE_URL=http://127.0.0.1:8080/api/v4
  TF_ACC=1
  ```

* **Useful HashiCorp documentation:**

  Refer to [HashiCorp's testing guide](https://www.terraform.io/docs/extend/testing/index.html)
  and [HashiCorp's testing best practices](https://www.terraform.io/docs/extend/best-practices/testing.html).
