#!/usr/bin/env sh

# NOTE: This script becomes obsolete when https://github.com/golang/go/issues/27348 is implemented.

set -e

go mod vendor
diff="$(git status --short ./vendor)"
test -z "$diff" || {
  echo "$diff"
  echo 'ERROR: Vendored Go modules are not checked in. Run "go mod vendor" and check in the changes.'
  exit 1
}
