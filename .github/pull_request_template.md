## PR Checklist

For a smooth review process, please run through this checklist before submitting a PR.

- [ ] Resource attributes match 1:1 the names and structure of the API resource in [the GitLab API documentation](https://docs.gitlab.com/ee/api/).
- [ ] Docs are updated with any new resources or attributes, including how to import the resource.
- [ ] New resources should have at minimum a basic test with three steps:
    - Create the resource
    - Update the attributes
    - Import the resource
- [ ] No new `//lintignore` comments that came from copied code. Linter rules are meant to be enforced on new code.
