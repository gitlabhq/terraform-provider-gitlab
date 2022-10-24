data "gitlab_group_subgroups" "subgroups" {
    group_id = "123456"
}

output "subgroups" {
    value = data.gitlab_group_subgroups.subgroups
}
