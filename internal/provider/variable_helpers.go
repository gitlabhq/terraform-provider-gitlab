package provider

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func augmentVariableClientError(d *schema.ResourceData, err error) diag.Diagnostics {
	// Masked values will commonly error due to their strict requirements, and the error message from the GitLab API is not very informative,
	// so we return a custom error message in this case.
	if d.Get("masked").(bool) && isInvalidValueError(err) {
		log.Printf("[ERROR] %v", err)
		return diag.Errorf("Invalid value for a masked variable. Check the masked variable requirements: https://docs.gitlab.com/ee/ci/variables/#masked-variable-requirements")
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func isInvalidValueError(err error) bool {
	var httpErr *gitlab.ErrorResponse
	return errors.As(err, &httpErr) &&
		httpErr.Response.StatusCode == http.StatusBadRequest &&
		strings.Contains(httpErr.Message, "value") &&
		strings.Contains(httpErr.Message, "invalid")
}
