resource "gitlab_group" "this" {
  name        = "example"
  path        = "example"
  description = "An example group"
}

resource "gitlab_project" "this" {
  name                   = "example"
  namespace_id           = gitlab_group.this.id
  initialize_with_readme = true
}

resource "gitlab_project_environment" "this" {
  project      = gitlab_project.this.id
  name         = "example"
  external_url = "www.example.com"
}

resource "gitlab_project_protected_environment" "this" {
  project     = gitlab_project.this.id
  environment = gitlab_project_environment.this.name

  deploy_access_levels {
    access_level = "developer"
  }
}

resource "gitlab_project_protected_environment" "this" {
  project     = gitlab_project.this.id
  environment = gitlab_project_environment.this.name

  deploy_access_levels {
    group_id = gitlab_group.test.id
  }
}

resource "gitlab_project_protected_environment" "this" {
  project     = gitlab_project.this.id
  environment = gitlab_project_environment.this.name

  deploy_access_levels {
    user_id = gitlab_user.test.id
  }

}
