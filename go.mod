module github.com/gitlabhq/terraform-provider-gitlab

go 1.16

require (
	github.com/hashicorp/go-retryablehttp v0.7.0
	github.com/hashicorp/terraform-plugin-sdk v1.16.0
	github.com/mitchellh/hashstructure v1.0.0
	github.com/xanzy/go-gitlab v0.50.0
)

replace github.com/xanzy/go-gitlab v0.46.0 => github.com/sirlatrom/go-gitlab v0.13.1-0.20210413075637-2b2ef8b0c50d
