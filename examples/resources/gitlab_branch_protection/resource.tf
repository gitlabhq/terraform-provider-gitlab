resource "gitlab_branch_protection" "BranchProtect" {
  project                      = "12345"
  branch                       = "BranchProtected"
  push_access_level            = "developer"
  merge_access_level           = "developer"
  allow_force_push             = true
  code_owner_approval_required = true
  allowed_to_push {
    user_id = 5
  }
  allowed_to_push {
    user_id = 521
  }
  allowed_to_merge {
    user_id = 15
  }
  allowed_to_merge {
    user_id = 37
  }
}

# Example using dynamic block
resource "gitlab_branch_protection" "main" {
  project            = "12345"
  branch             = "main"
  push_access_level  = "maintainer"
  merge_access_level = "maintainer"

  dynamic "allowed_to_push" {
    for_each = [50, 55, 60]
    content {
      user_id = allowed_to_push.value
    }
  }
}
