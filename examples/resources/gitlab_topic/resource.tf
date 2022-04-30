resource "gitlab_topic" "functional_programming" {
  name        = "functional-programming"
  title       = "Functional Programming"
  description = "In computer science, functional programming is a programming paradigm where programs are constructed by applying and composing functions."
  avatar      = "${path.module}/avatar.png"
  avatar_hash = filesha256("${path.module}/avatar.png")
}
