package gitlab

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/xanzy/go-gitlab"
	"regexp"
)

// copied from ../github/util.go
func validateValueFunc(values []string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (we []string, errors []error) {
		value := v.(string)
		valid := false
		for _, role := range values {
			if value == role {
				valid = true
				break
			}
		}

		if !valid {
			errors = append(errors, fmt.Errorf("%s is an invalid value for argument %s", value, k))
		}
		return
	}
}

func stringToVisibilityLevel(s string) *gitlab.VisibilityValue {
	lookup := map[string]gitlab.VisibilityValue{
		"private":  gitlab.PrivateVisibility,
		"internal": gitlab.InternalVisibility,
		"public":   gitlab.PublicVisibility,
	}

	value, ok := lookup[s]
	if !ok {
		return nil
	}
	return &value
}

func StringIsGitlabVariableName() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (s []string, es []error) {
		value, ok := v.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}
		if len(value) < 1 || len(value) > 255 {
			es = append(es, fmt.Errorf("expected length of %s to be in the range (%d - %d), got %s", k, 1, 255, v))
		}

		match, _ := regexp.MatchString("[a-zA-Z0-9_]+", value)
		if !match {
			es = append(es, fmt.Errorf("%s is an invalid value for argument %s. Only A-Z, a-z, 0-9, and _ are allowed", value, k))
		}
		return
	}
}
