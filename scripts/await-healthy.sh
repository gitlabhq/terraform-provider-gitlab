#!/usr/bin/env sh

echo 'Waiting for GitLab container to become healthy'

until test -n "$(docker ps --quiet --filter label=terraform-provider-gitlab/owned --filter health=healthy)"; do
  printf '.'
  sleep 5
done

echo 'GitLab is healthy'
