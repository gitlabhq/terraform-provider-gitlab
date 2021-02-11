module github.com/gitlabhq/terraform-provider-gitlab

go 1.16

require (
	github.com/bflad/tfproviderlint v0.27.0
	github.com/hashicorp/go-retryablehttp v0.7.0
	github.com/hashicorp/terraform-plugin-sdk v1.16.1
	github.com/mitchellh/hashstructure v1.0.0
	github.com/xanzy/go-gitlab v0.50.0
)

replace github.com/xanzy/go-gitlab v0.42.0 => github.com/sirlatrom/go-gitlab v0.13.1-0.20210210175321-27539466b2c8
