data "gitlab_current_user" "me" {
}

// https://github.com/terraform-providers/terraform-provider-gitlab/issues/321
resource "gitlab_application" "grafana" {
  // create the application if user has admin role, 
  // otherwise Gitlab API rejects requests with 403 http status code
  count = data.gitlab_current_user.me.is_admin ? 1 : 0

  name = "grafana"
  redirect_uri = [
    "http://localhost:8080/auth/gitlab/callback",
    "http://localhost:8080/auth/gitlab/callback2"
  ]
  scopes = ["openid", "email", "profile"]
}

output "me" {
  value = data.gitlab_current_user.me
}

output "grafana_app_secret" {
  value = {
    app_id     = gitlab_application.grafana[0].application_id
    app_secret = gitlab_application.grafana[0].secret
  }
  sensitive = true
}
