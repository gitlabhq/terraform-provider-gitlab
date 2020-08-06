module github.com/terraform-providers/terraform-provider-gitlab

replace github.com/xanzy/go-gitlab v0.32.1 => github.com/sfang97/go-gitlab v0.33.1-0.20200806213917-abe9c5791dd4

require (
	github.com/hashicorp/terraform-plugin-sdk v1.13.1
	github.com/mitchellh/hashstructure v1.0.0
	github.com/xanzy/go-gitlab v0.32.1
)

go 1.14
