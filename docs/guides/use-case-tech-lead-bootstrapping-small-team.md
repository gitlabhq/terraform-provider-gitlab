---
page_title: "A Tech Lead Bootstrapping a Small Team"
subcategory: "Use Cases"
---

Imagine you are a tech lead responsible for a small team, and you want to get
your team bootstrapped with their own group, a wiki, and a couple of projects: 
one to hold the full-stack application that you've been working on, and another 
for a user facing documentation website. You want to make sure that code quality 
is verified by yourself and at least one additional team member. You also want 
to set up your own GitLab CI/CD runner to run your automation jobs for your 
full-stack mono-repo as well as your documentation website.

You've done some research and learned that IaC (Infrastructure as Code) is all 
the rage. You also know that OpenTofu and Terraform are a great technologies for 
cloud resources. In addition, you know that GitLab, your SDLC tool of choice, 
has its own Terraform provider that will enable you to realize all of your IaC 
dreams! You still have a problem though: how in the world can you take advantage 
of this?

Have no fear! This guide will walk you through the process of starting from a 
fresh installation of OpenTofu or Terraform and building out your infrastructure 
as code solution.

# Configure GitLab provider

Create a new directory on your computer to store your IaC code. From that 
directory, create a file named `main.tf`, and add the following lines:

```terraform
terraform {
  required_providers {
    gitlab = {
      source = "gitlabhq/gitlab"
    }
  }
}

variable "admin_token" {
  description = "Owner/Maintainer PAT token with the api scope applied."
  type        = string
}

provider "gitlab" {
  token    = var.admin_token
  base_url = "https://gitlab.com" # change this if you are on a self-hosted GitLab instance.
}
```

The `terraform` block tells OpenTofu or Terraform where to download the provider 
for all GitLab resources.

The `provider` block configures the provider to use an externally provided 
personal access token to authenticate with GitLab when performing any
configuration.

# Create a group and projects

Next, let's create a GitLab group for your team's code and wiki. Groups allow
you to supply a multi-level hierarchy to your code assets. A root, or top-level
group, is typically provided by an organization so they have policy-level
controls over all sub-groups and projects. It's a best practice to limit the
number of top-level groups, and to subdivide your groups into team or functional
areas. To facilitate this, create a group for the team by adding these lines to
your `main.tf` file, modifying it to fit your needs:

```terraform
resource "gitlab_group" "my_team" {
  parent_id         = 1337           # change to your top-level group ID number
  name              = "Awesome Team" # friendly group name
  path              = "awesome-team" # path that will be a part of clone URIs
  description       = "The Awesome Team provides awesome tech that makes our company shine!"
  wiki_access_level = "private"      # make the Wiki only viewable by group members
}
```

Now that we have a group, we can add some team members to it. GitLab provides
access levels that you can apply to team members to give them different
permissions. As the team lead, we'll give you _maintainer_ access. For your team
members, we'll give them _developer_ access. Because it's a pain to figure out
the internal ID of each team member, we'll take advantage of a data source to
supply user handles, instead of using ID numbers.

Add this code to the bottom of your `main.tf`, modifying it to fit your needs:

```terraform
data "gitlab_user" "team_lead" {
  username = "Delaney"
}

data "gitlab_user" "team_members" {
  for_each = toset(["Sasha", "Priyanka", "Simone"])
  username = each.value
}

resource "gitlab_group_membership" "team_lead" {
  group_id     = gitlab_group.my_team.id
  access_level = "maintainer"
  user_id      = data.gitlab_user.team_lead.id
}

resource "gitlab_group_membership" "team_members" {
  for_each     = data.gitlab_user.team_members
  group_id     = gitlab_group.my_team.id
  access_level = "developer"
  user_id      = each.value.id
}
```

Because of how membership rights work in GitLab, any projects we create under
this group will give these users the level of access defined at the group level.
So let's create one project for your full-stack app, and another for user-facing
documentation.

Add this code to the bottom of your `main.tf`, modifying it to fit your needs:

```terraform
resource "gitlab_project" "app" {
  namespace_id = gitlab_group.my_team.id
  name         = "Fullstack App"
  path         = "app"
  description  = "Our fullstack app which will deliver value fast!"
  wiki_enabled = false
}

resource "gitlab_project" "docs" {
  namespace_id = gitlab_group.my_team.id
  name         = "Documentation"
  path         = "docs"
  description  = "User facing documentation website."
  wiki_enabled = false
}
```

Now, let's add some approval rules. These force any MRs created in these
projects to be approved by you (the tech lead) and at least one other team
member before they can be merged. You must have at least GitLab Premium for your
top-level group to complete this step. If you don't have GitLab Premium, you can
skip ahead.

Add this code to the bottom of your `main.tf`, modifying it to fit your needs:

```terraform
resource "gitlab_project_approval_rule" "team_app_maintainers" {
  project            = gitlab_project.app.id
  name               = "maintainers"
  approvals_required = 1
  user_ids           = [data.gitlab_user.team_lead.id]
}

resource "gitlab_project_approval_rule" "team_app_members" {
  project            = gitlab_project.app.id
  name               = "members"
  approvals_required = 1
  user_ids           = [for user in data.gitlab_user.team_members : user.id]
}

resource "gitlab_project_approval_rule" "team_docs_maintainers" {
  project            = gitlab_project.docs.id
  name               = "maintainers"
  approvals_required = 1
  user_ids           = [data.gitlab_user.team_lead.id]
}

resource "gitlab_project_approval_rule" "team_docs_members" {
  project            = gitlab_project.docs.id
  name               = "members"
  approvals_required = 1
  user_ids           = [for user in data.gitlab_user.team_members : user.id]
}
```

# Configure CI/CD runner

With this in place, we have one item left. You want to be able to automatically
run tests, build your product, package it and ship it to customers. For that,
you need a GitLab runner! With the current runner registration workflow, you
must create a runner instance on your GitLab group or project in order to
configure basic settings, and to get a registration token that you can use with
your deployed runners. We are going to create a group runner so that it can be
shared with your full-stack application and user documentation projects.

Add this code to the bottom of your `main.tf`, modifying it to fit your needs:

```terraform
resource "gitlab_user_runner" "linux" {
  group_id    = gitlab_group.my_team.id
  description = "Team Linux Job Runner"
  runner_type = "group_type"
  untagged    = true # you can use `tag_list` instead if you want user's to opt-in to using this runner.
}

output "registration_token" {
  description = "Registration token to to use with your runner installation."
  value       = gitlab_user_runner.linux.token
  sensitive   = true
}
```

With all of this configuration in place, you should be ready to rock. You can
use the following commands to initialize your OpenTofu or Terraform root module,
review the changes, apply them, and then retrieve the registration token for
your runner installation:

```shell
tofu init
tofu plan -out plan.out
tofu apply plan.out
tofu output registration_token
```

Note that this example uses OpenTofu's `tofu` command, but you can easily switch
it out with `terraform` instead. Terraform and OpenTofu are backwards and
forwards compatible.
