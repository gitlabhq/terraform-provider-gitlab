data "gitlab_user" "example" {
  username = "myuser"
}

# Example using `for_each`
data "gitlab_user" "example-two" {
  for_each = toset(["user1", "user2", "user3"])

  username = each.value
}
