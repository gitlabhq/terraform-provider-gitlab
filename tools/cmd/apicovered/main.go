package main

import (
	"os"

	"github.com/gitlabhq/terraform-provider-gitlab/tools/passes/apicovered"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	apicovered.Output = os.Stdout
	singlechecker.Main(apicovered.Analyzer)
}
