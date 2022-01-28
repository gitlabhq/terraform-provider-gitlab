package gitlab

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func renderValueListForDocs(values []string) string {
	inlineCodeValues := make([]string, 0, len(values))
	for _, v := range values {
		inlineCodeValues = append(inlineCodeValues, fmt.Sprintf("`%s`", v))
	}
	return strings.Join(inlineCodeValues, ", ")
}

var validateDateFunc = func(v interface{}, k string) (we []string, errors []error) {
	value := v.(string)
	//add zero hours and let time figure out correctness
	_, e := time.Parse(time.RFC3339, value+"T00:00:00Z")
	if e != nil {
		errors = append(errors, fmt.Errorf("%s is not valid for format YYYY-MM-DD", value))
	}
	return
}

var validateURLFunc = func(v interface{}, k string) (s []string, errors []error) {
	value := v.(string)
	url, err := url.Parse(value)

	if err != nil || url.Host == "" || url.Scheme == "" {
		errors = append(errors, fmt.Errorf("%s is not a valid URL", value))
		return
	}

	return
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

func stringToProjectCreationLevel(s string) *gitlab.ProjectCreationLevelValue {
	lookup := map[string]gitlab.ProjectCreationLevelValue{
		"noone":      gitlab.NoOneProjectCreation,
		"maintainer": gitlab.MaintainerProjectCreation,
		"developer":  gitlab.DeveloperProjectCreation,
	}

	value, ok := lookup[s]
	if !ok {
		return nil
	}
	return &value
}

func stringToSubGroupCreationLevel(s string) *gitlab.SubGroupCreationLevelValue {
	lookup := map[string]gitlab.SubGroupCreationLevelValue{
		"owner":      gitlab.OwnerSubGroupCreationLevelValue,
		"maintainer": gitlab.MaintainerSubGroupCreationLevelValue,
	}

	value, ok := lookup[s]
	if !ok {
		return nil
	}
	return &value
}

func stringToVariableType(s string) *gitlab.VariableTypeValue {
	lookup := map[string]gitlab.VariableTypeValue{
		"env_var": gitlab.EnvVariableType,
		"file":    gitlab.FileVariableType,
	}

	value, ok := lookup[s]
	if !ok {
		return nil
	}
	return &value
}

func stringToMergeMethod(s string) *gitlab.MergeMethodValue {
	lookup := map[string]gitlab.MergeMethodValue{
		"merge":        gitlab.NoFastForwardMerge,
		"ff":           gitlab.FastForwardMerge,
		"rebase_merge": gitlab.RebaseMerge,
	}

	value, ok := lookup[s]
	if !ok {
		return nil
	}
	return &value
}

func stringToSquashOptionValue(s string) *gitlab.SquashOptionValue {
	lookup := map[string]gitlab.SquashOptionValue{
		"never":       gitlab.SquashOptionNever,
		"always":      gitlab.SquashOptionAlways,
		"default_on":  gitlab.SquashOptionDefaultOn,
		"default_off": gitlab.SquashOptionDefaultOff,
	}

	value, ok := lookup[s]
	if !ok {
		return nil
	}
	return &value
}

func stringToAccessControlValue(s string) *gitlab.AccessControlValue {
	lookup := map[string]gitlab.AccessControlValue{
		"disabled": gitlab.DisabledAccessControl,
		"enabled":  gitlab.EnabledAccessControl,
		"private":  gitlab.PrivateAccessControl,
		"public":   gitlab.PublicAccessControl,
	}

	value, ok := lookup[s]
	if !ok {
		return nil
	}
	return &value
}

// lintignore: V011 // TODO: Resolve this tfproviderlint issue
var StringIsGitlabVariableName = func(v interface{}, k string) (s []string, es []error) {
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

var StringIsGitlabVariableType = func(v interface{}, k string) (s []string, es []error) {
	value, ok := v.(string)
	if !ok {
		es = append(es, fmt.Errorf("expected type of %s to be string", k))
		return
	}
	variableType := stringToVariableType(value)
	if variableType == nil {
		es = append(es, fmt.Errorf("expected variable_type to be \"env_var\" or \"file\""))
	}
	return
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

var tagProtectionAccessLevelID = map[string]gitlab.AccessLevelValue{
	"no one":     gitlab.NoPermissions,
	"developer":  gitlab.DeveloperPermissions,
	"maintainer": gitlab.MaintainerPermissions,
}

var tagProtectionAccessLevelNames = map[gitlab.AccessLevelValue]string{
	gitlab.NoPermissions:         "no one",
	gitlab.DeveloperPermissions:  "developer",
	gitlab.MaintainerPermissions: "maintainer",
}

func stringSetToStringSlice(stringSet *schema.Set) *[]string {
	ret := []string{}
	if stringSet == nil {
		return &ret
	}
	for _, envVal := range stringSet.List() {
		ret = append(ret, envVal.(string))
	}
	return &ret
}

// isGitLabVersionLessThan is a SkipFunc that returns true if the provided version is lower then
// the current version of GitLab. It only checks the major and minor version numbers, not the patch.
func isGitLabVersionLessThan(client *gitlab.Client, version string) func() (bool, error) {
	return func() (bool, error) {
		isAtLeast, err := isGitLabVersionAtLeast(client, version)()
		return !isAtLeast, err
	}
}

// isGitLabVersionAtLeast is a SkipFunc that checks that the version of GitLab is at least the
// provided wantVersion. It only checks the major and minor version numbers, not the patch.
func isGitLabVersionAtLeast(client *gitlab.Client, wantVersion string) func() (bool, error) {
	return func() (bool, error) {
		wantMajor, wantMinor, err := parseVersionMajorMinor(wantVersion)
		if err != nil {
			return false, fmt.Errorf("failed to parse wanted version %q: %w", wantVersion, err)
		}

		actualVersion, _, err := client.Version.GetVersion()
		if err != nil {
			return false, err
		}

		actualMajor, actualMinor, err := parseVersionMajorMinor(actualVersion.Version)
		if err != nil {
			return false, fmt.Errorf("failed to parse actual version %q: %w", actualVersion.Version, err)
		}

		if actualMajor == wantMajor {
			return actualMinor >= wantMinor, nil
		}

		return actualMajor > wantMajor, nil
	}
}

func parseVersionMajorMinor(version string) (int, int, error) {
	parts := strings.SplitN(version, ".", 3)
	if len(parts) < 2 {
		return 0, 0, fmt.Errorf("need at least 2 parts (was %d)", len(parts))
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}

	return major, minor, nil
}

func is404(err error) bool {
	if errResponse, ok := err.(*gitlab.ErrorResponse); ok &&
		errResponse.Response != nil &&
		errResponse.Response.StatusCode == 404 {
		return true
	}
	return false
}
