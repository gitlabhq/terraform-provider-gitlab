package provider

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

// extractIIDFromGlobalID extracts the internal model ID from a global GraphQL ID.
//
// e.g. 'gid://gitlab/User/1' -> 1 or 'gid://gitlab/Project/42' -> 42
//
// see https://docs.gitlab.com/ee/development/api_graphql_styleguide.html#global-ids
func extractIIDFromGlobalID(globalID string) (int, error) {
	parts := strings.Split(globalID, "/")
	iid, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return 0, fmt.Errorf("unable to extract iid from global id %q. Was looking for an integer after the last slash (/).", globalID)
	}
	return iid, nil
}

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

func stringListToStringSlice(stringList []interface{}) *[]string {
	ret := []string{}
	if stringList == nil {
		return &ret
	}
	for _, v := range stringList {
		ret = append(ret, fmt.Sprint(v))
	}
	return &ret
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

func intSetToIntSlice(intSet *schema.Set) *[]int {
	ret := []int{}
	if intSet == nil {
		return &ret
	}
	for _, envVal := range intSet.List() {
		ret = append(ret, envVal.(int))
	}
	return &ret
}

// isGitLabVersionLessThan is a SkipFunc that returns true if the provided version is lower then
// the current version of GitLab. It only checks the major and minor version numbers, not the patch.
func isGitLabVersionLessThan(ctx context.Context, client *gitlab.Client, version string) func() (bool, error) {
	return func() (bool, error) {
		isAtLeast, err := isGitLabVersionAtLeast(ctx, client, version)()
		return !isAtLeast, err
	}
}

// isGitLabVersionAtLeast is a SkipFunc that checks that the version of GitLab is at least the
// provided wantVersion. It only checks the major and minor version numbers, not the patch.
func isGitLabVersionAtLeast(ctx context.Context, client *gitlab.Client, wantVersion string) func() (bool, error) {
	return func() (bool, error) {
		wantMajor, wantMinor, err := parseVersionMajorMinor(wantVersion)
		if err != nil {
			return false, fmt.Errorf("failed to parse wanted version %q: %w", wantVersion, err)
		}

		actualVersion, _, err := client.Version.GetVersion(gitlab.WithContext(ctx))
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

func isCurrentUserAdmin(ctx context.Context, client *gitlab.Client) (bool, error) {
	currentUser, _, err := client.Users.CurrentUser(gitlab.WithContext(ctx))
	if err != nil {
		return false, err
	}

	return currentUser.IsAdmin, nil
}

// ISO 8601 date format
const iso8601 = "2006-01-02"

// isISO8601 validates if the given value is a ISO8601 compatible date in the YYYY-MM-DD format.
func isISO6801Date(i interface{}, p cty.Path) diag.Diagnostics {
	v := i.(string)

	if _, err := time.Parse(iso8601, v); err != nil {
		return diag.Errorf("expected %q to be a valid YYYY-MM-DD date, got %q: %+v", p, i, err)
	}

	return nil
}

func parseISO8601Date(v string) (*gitlab.ISOTime, error) {
	iso8601Date, err := time.Parse(iso8601, v)
	if err != nil {
		return nil, fmt.Errorf("expected %q to be a valid YYYY-MM-DD date", v)
	}

	x := gitlab.ISOTime(iso8601Date)
	return &x, nil
}

// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func constructSchema(schemas ...map[string]*schema.Schema) map[string]*schema.Schema {
	schema := make(map[string]*schema.Schema)
	for _, s := range schemas {
		for k, v := range s {
			schema[k] = v
		}
	}
	return schema
}

func attributeNamesFromSchema(schema map[string]*schema.Schema) []string {
	names := make([]string, 0, len(schema))
	for name := range schema {
		names = append(names, name)
	}
	return names
}

// datasourceSchemaFromResourceSchema is a recursive func that
// converts an existing Resource schema to a Datasource schema.
// All schema elements are copied, but certain attributes are ignored or changed:
// - all attributes have Computed = true
// - all attributes have ForceNew, Required = false
// - Validation funcs and attributes (e.g. MaxItems) are not copied
// Adapted from https://github.com/hashicorp/terraform-provider-google/blob/1a72f93a8dcf6f1e59d5f25aefcb6d794a116bf5/google/datasource_helpers.go#L13
func datasourceSchemaFromResourceSchema(rs map[string]*schema.Schema, arguments []string, optionalArguments []string) map[string]*schema.Schema {
	ds := make(map[string]*schema.Schema, len(rs))
	for k, v := range rs {
		dv := &schema.Schema{
			ForceNew:    false,
			Description: v.Description,
			Type:        v.Type,
		}
		if contains(arguments, k) {
			dv.Computed = false
			dv.Required = true
		} else {
			dv.Computed = true
			dv.Required = false
		}

		if contains(optionalArguments, k) {
			dv.Optional = true
		}

		switch v.Type {
		case schema.TypeSet:
			dv.Set = v.Set
			fallthrough
		case schema.TypeList:
			// List & Set types are generally used for 2 cases:
			// - a list/set of simple primitive values (e.g. list of strings)
			// - a sub resource
			if elem, ok := v.Elem.(*schema.Resource); ok {
				// handle the case where the Element is a sub-resource
				dv.Elem = &schema.Resource{
					Schema: datasourceSchemaFromResourceSchema(elem.Schema, nil, nil),
				}
			} else {
				// handle simple primitive case
				dv.Elem = v.Elem
			}

		default:
			// Elem of all other types are copied as-is
			dv.Elem = v.Elem

		}
		ds[k] = dv

	}
	return ds
}

func setStateMapInResourceData(stateMap map[string]interface{}, d *schema.ResourceData) error {
	for k, v := range stateMap {
		// lintignore: R001 // for convenience sake, to reduce maintenance burden we are ok not having literals here.
		if err := d.Set(k, v); err != nil {
			return fmt.Errorf("failed to set state for %q to %v: %w", k, v, err)
		}
	}

	return nil
}

// lock can be used to lock, but make it `context.Context` aware.
// e.g. it'll respect cancelling and timeouts.
type lock chan struct{}

func newLock() lock {
	return make(lock, 1)

}

func (c lock) lock(ctx context.Context) error {
	select {
	case c <- struct{}{}:
		// lock acquired
		return nil
	case <-ctx.Done():
		// Timeout
		return ctx.Err()
	}
}

func (c lock) unlock() {
	<-c
}
