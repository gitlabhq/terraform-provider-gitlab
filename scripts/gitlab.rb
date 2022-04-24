registry_external_url 'http://127.0.0.1:5050'
registry['enable']                    = true
registry_nginx['ssl_certificate']     = "/etc/gitlab/ssl/gitlab-registry.pem"
registry_nginx['ssl_certificate_key'] = "/etc/gitlab/ssl/gitlab-registry.key"
