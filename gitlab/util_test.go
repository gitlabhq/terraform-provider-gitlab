package gitlab

import (
	"testing"

	gitlab "github.com/xanzy/go-gitlab"
)

func TestGitlab_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "invalid",
			ErrCount: 1,
		},
		{
			Value:    "valid_one",
			ErrCount: 0,
		},
		{
			Value:    "valid_two",
			ErrCount: 0,
		},
	}

	validationFunc := validateValueFunc([]string{"valid_one", "valid_two"})

	for _, tc := range cases {
		_, errors := validationFunc(tc.Value, "test_arg")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected 1 validation error")
		}
	}
}

func TestGitlab_visbilityHelpers(t *testing.T) {
	cases := []struct {
		String string
		Level  gitlab.VisibilityValue
	}{
		{
			String: "private",
			Level:  gitlab.PrivateVisibility,
		},
		{
			String: "public",
			Level:  gitlab.PublicVisibility,
		},
	}

	for _, tc := range cases {
		level := stringToVisibilityLevel(tc.String)
		if level == nil || *level != tc.Level {
			t.Fatalf("got %v expected %v", level, tc.Level)
		}

		sv := string(tc.Level)
		if sv == "" || sv != tc.String {
			t.Fatalf("got %v expected %v", sv, tc.String)
		}
	}
}

func TestValidateURLFunc(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "invalid_url",
			ErrCount: 1,
		},
		{
			Value:    "invalid_url.com",
			ErrCount: 1,
		},
		{
			Value:    "/relative/path",
			ErrCount: 1,
		},
		{
			Value:    "https://valid_url.com",
			ErrCount: 0,
		},
		{
			Value:    "http://www.valid_url.com",
			ErrCount: 0,
		},
	}

	validationFunc := validateURLFunc()

	for _, tc := range cases {
		_, errors := validationFunc(tc.Value, "test_arg")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected 1 validation error")
		}
	}
}
