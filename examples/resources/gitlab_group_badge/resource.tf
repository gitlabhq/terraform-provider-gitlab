resource "gitlab_group" "foo" {
  name = "foo-group"
}

resource "gitlab_group_badge" "example" {
  group     = gitlab_group.foo.id
  link_url  = "https://example.com/badge-123"
  image_url = "https://example.com/badge-123.svg"
}
