## Description

<!-- Which issue/s does this PR close? Is there any more context you can give the reviewer? -->

## PR Checklist

<!-- For a smooth review process, please run through this checklist before submitting a PR. -->

- [ ] Resource attributes match 1:1 the names and structure of the API resource in [the GitLab API documentation](https://docs.gitlab.com/ee/api/).
- [ ] [Examples](/examples) are updated with:
    - A \*.tf file for the resource/s with at least one usage example
    - A \*.sh file for the resource/s with an import example (if applicable)
- [ ] New resources have at minimum a basic test with three steps:
    - Create the resource
    - Update the attributes
    - Import the resource
- [ ] No new `//lintignore` comments that came from copied code. Linter rules are meant to be enforced on new code.
