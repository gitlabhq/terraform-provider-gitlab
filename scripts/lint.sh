#!/usr/bin/env bash

docker run --rm -v "$(pwd):/workspace" -w /workspace golangci/golangci-lint:v1.41.1 golangci-lint run
docker run --rm -v "$(pwd):/workspace" -w /workspace avtodev/markdown-lint:v1 --config .markdownlint.yml docs
