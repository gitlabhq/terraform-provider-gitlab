package tools

//go:generate go install github.com/bflad/tfproviderlint/cmd/tfproviderlintx
//go:generate go install github.com/golangci/golangci-lint/cmd/golangci-lint
//go:generate go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
//go:generate go install mvdan.cc/sh/v3/cmd/shfmt

import (
	_ "github.com/bflad/tfproviderlint/cmd/tfproviderlintx"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
	_ "mvdan.cc/sh/v3/cmd/shfmt"
)
