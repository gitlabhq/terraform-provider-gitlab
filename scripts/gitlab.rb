pages_external_url 'http://127.0.0.1:5051'
pages_nginx['redirect_http_to_https'] = false
pages_nginx['ssl_certificate'] = "/etc/gitlab/ssl/gitlab-registry.pem"
pages_nginx['ssl_certificate_key'] = "/etc/gitlab/ssl/gitlab-registry.key"

registry_external_url 'http://127.0.0.1:5050'
registry['enable']                    = true
registry_nginx['ssl_certificate']     = "/etc/gitlab/ssl/gitlab-registry.pem"
registry_nginx['ssl_certificate_key'] = "/etc/gitlab/ssl/gitlab-registry.key"
