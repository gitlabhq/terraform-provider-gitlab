#!/usr/bin/env sh

# This script is run during the CI workflow in order to run additional custom checks and validations
# that do not fit into an existing lint framework.
#
# This script uses GitHub workflow commands in order to add contextual error messages:
#   https://docs.github.com/en/actions/using-workflows/workflow-commands-for-github-actions#setting-an-error-message
#
# Usage:
#   ./scripts/lint-custom.sh

has_failure=''

pattern="^\*\*Upstream API\*\*: \[.*\](.*)$"
files_without_match="$(grep --recursive --files-without-match "$pattern" ./docs/resources ./docs/data-sources)"
test -z "$files_without_match" || {
    has_failure='true'
    echo "$files_without_match" | xargs -I{} echo "::error file={}::Generated documentation file {} should reference an upstream GitLab API. Expected file to match regex \"$pattern\"."
}

# ^^^ Add new checks above this line ^^^

test -z "$has_failure"
