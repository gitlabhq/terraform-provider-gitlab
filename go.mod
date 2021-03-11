module github.com/gitlabhq/terraform-provider-gitlab

go 1.16

replace github.com/xanzy/go-gitlab => github.com/randomswdev/go-gitlab v0.46.1-0.20210311212325-5c6b3b46bea4

require (
	github.com/bflad/tfproviderlint v0.27.0
	github.com/hashicorp/go-retryablehttp v0.7.0
	github.com/hashicorp/terraform-plugin-sdk v1.16.1
	github.com/mitchellh/hashstructure v1.0.0
	github.com/onsi/gomega v1.14.0
	github.com/xanzy/go-gitlab v0.50.0
)
