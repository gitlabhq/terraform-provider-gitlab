<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Terraform Provider for Gitlab
=============================

- [Documentation](https://www.terraform.io/docs/providers/gitlab/index.html)
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)
- Build status:
  - ![Unit Tests](https://github.com/terraform-providers/terraform-provider-gitlab/workflows/Unit%20Tests/badge.svg)
  - ![Acceptance Tests](https://github.com/terraform-providers/terraform-provider-gitlab/workflows/Acceptance%20Tests/badge.svg)
  - ![Website Build](https://github.com/terraform-providers/terraform-provider-gitlab/workflows/Website%20Build/badge.svg)

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.12 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-gitlab`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-gitlab
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-gitlab
$ make build
```

Using the provider
----------------------
## Fill in for each provider

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-gitlab
...
```

### Running tests
The Terraform Provider only has acceptance tests, these can run against a gitlab instance where you have a token with administrator permissions (likely not gitlab.com). 
There is excellent documentation on [how to run gitlab from docker at gitlab.com](https://docs.gitlab.com/omnibus/docker/)

In order to run the full suite of acceptance tests, export the environment variables: 

- `GITLAB_TOKEN` //token for account with admin priviliges
- `GITLAB_BASE_URL` //URL with api part e.g. `http://localhost:8929/api/v4/`

and run `make testacc`.

```sh
$ make testacc
```
