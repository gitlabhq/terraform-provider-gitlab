package provider

import (
	"testing"

	gitlab "github.com/xanzy/go-gitlab"
)

func TestGitlab_extractIIDFromGlobalID(t *testing.T) {
	cases := []struct {
		GlobalID string
		IID      int
	}{
		{
			GlobalID: "gid://gitlab/User/1",
			IID:      1,
		},
		{
			GlobalID: "gid://gitlab/Namespaces::UserNamespace/1000",
			IID:      1000,
		},
	}

	for _, tc := range cases {
		iid, err := extractIIDFromGlobalID(tc.GlobalID)
		if err != nil {
			t.Fatalf("expected valid global id, got %q: %v", tc.GlobalID, err)
		}

		if iid != tc.IID {
			t.Fatalf("got %v expected %v", iid, tc.IID)
		}
	}
}

func TestGitlab_extractIIDFromGlobalID_invalidGlobalID(t *testing.T) {
	cases := []struct {
		GlobalID string
	}{
		{
			GlobalID: "",
		},
		{
			GlobalID: "gid://gitlab/User/",
		},
		{
			GlobalID: "gid://gitlab/Namespaces::UserNamespace",
		},
	}

	for _, tc := range cases {
		iid, err := extractIIDFromGlobalID(tc.GlobalID)
		if err == nil {
			t.Fatalf("expected invalid global id, got id %q instead from global id %q", iid, tc.GlobalID)
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

	for _, tc := range cases {
		_, errors := validateURLFunc(tc.Value, "test_arg")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected 1 validation error")
		}
	}
}
