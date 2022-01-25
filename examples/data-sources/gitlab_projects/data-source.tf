# List projects within a group tree
data "gitlab_group" "mygroup" {
  full_path = "mygroup"
}

data "gitlab_projects" "group_projects" {
  group_id          = data.gitlab_group.mygroup.id
  order_by          = "name"
  include_subgroups = true
  with_shared       = false
}

# List projects using the search syntax
data "gitlab_projects" "projects" {
  search              = "postgresql"
  visibility          = "private"
}
