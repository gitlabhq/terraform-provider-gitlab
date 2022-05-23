resource "gitlab_project_environment" "this" {
  project      = 123
  name         = "example"
  external_url = "www.example.com"
}

# Example with access level
resource "gitlab_project_protected_environment" "example_with_access_level" {
  project                 = gitlab_project_environment.this.project
  required_approval_count = 1
  environment             = gitlab_project_environment.this.name

  deploy_access_levels {
    access_level = "developer"
  }
}

# Example with group
resource "gitlab_project_protected_environment" "example_with_group" {
  project     = gitlab_project_environment.this.project
  environment = gitlab_project_environment.this.name

  deploy_access_levels {
    group_id = 456
  }
}

# Example with user
resource "gitlab_project_protected_environment" "example_with_user" {
  project     = gitlab_project_environment.this.project
  environment = gitlab_project_environment.this.name

  deploy_access_levels {
    user_id = 789
  }
}

# Example with multiple access levels
resource "gitlab_project_protected_environment" "example_with_multiple" {
  project                 = gitlab_project_environment.this.project
  required_approval_count = 2
  environment             = gitlab_project_environment.this.name

  deploy_access_levels {
    access_level = "developer"
  }

  deploy_access_levels {
    group_id = 456
  }

  deploy_access_levels {
    user_id = 789
  }
}
