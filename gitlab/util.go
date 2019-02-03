package gitlab

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/xanzy/go-gitlab"
	"regexp"
	"strings"
	"time"
)

var accessLevelNameToValue = map[string]gitlab.AccessLevelValue{
	"guest":     gitlab.GuestPermissions,
	"reporter":  gitlab.ReporterPermissions,
	"developer": gitlab.DeveloperPermissions,
	"master":    gitlab.MasterPermissions,
	"owner":     gitlab.OwnerPermission,
}

var accessLevelValueToName = map[gitlab.AccessLevelValue]string{
	gitlab.GuestPermissions:     "guest",
	gitlab.ReporterPermissions:  "reporter",
	gitlab.DeveloperPermissions: "developer",
	gitlab.MasterPermissions:    "master",
	gitlab.OwnerPermission:      "owner",
}

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
			errors = append(errors, fmt.Errorf("%s is an invalid value for argument %s acceptable values are: %v", value, k, values))
		}
		return
	}
}

func validateDateFunc() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (we []string, errors []error) {
		value := v.(string)
		//add zero hours and let time figure out correctness
		_, e := time.Parse(time.RFC3339, value+"T00:00:00Z")
		if e != nil {
			errors = append(errors, fmt.Errorf("%s is not valid for format YYYY-MM-DD", value))
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

// return the pieces of id `a:b` as a, b
func parseTwoPartID(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("Unexpected ID format (%q). Expected project:key", id)
	}

	return parts[0], parts[1], nil
}

// format the strings into an id `a:b`
func buildTwoPartID(a, b *string) string {
	return fmt.Sprintf("%s:%s", *a, *b)
}

var accessLevelID = map[string]gitlab.AccessLevelValue{
	"guest":     gitlab.GuestPermissions,
	"reporter":  gitlab.ReporterPermissions,
	"developer": gitlab.DeveloperPermissions,
	"master":    gitlab.MasterPermissions,
	"owner":     gitlab.OwnerPermission,
}

var accessLevel = map[gitlab.AccessLevelValue]string{
	gitlab.GuestPermissions:     "guest",
	gitlab.ReporterPermissions:  "reporter",
	gitlab.DeveloperPermissions: "developer",
	gitlab.MasterPermissions:    "master",
	gitlab.OwnerPermission:      "owner",
}
