#!/bin/bash -e
test "$MAKE_TARGET" == "testacc" || { echo "not starting gitlab!"; exit 0; }
echo "Starting gitlab container..."
docker run -d --rm --name gitlab \
  -e GITLAB_ROOT_PASSWORD=adminadmin \
  -p 127.0.0.1:8080:80 \
  gitlab/gitlab-ce

echo -n "Waiting for gitlab to be ready "
i=1
until wget -t 1 127.0.0.1:8080 -O /dev/null -q
do
    sleep 1
    echo -n .
    if [[ $((i%3)) == 0 ]]; then echo -n ' '; fi
    (( i++ ))
done

echo
echo "Creating access token"
(
  echo -n 'terraform_token = PersonalAccessToken.create('
  echo -n 'user_id: 1, '
  echo -n 'scopes: [:api, :read_user], '
  echo -n 'name: :terraform);'
  echo -n "terraform_token.set_token('${GITLAB_TOKEN:=ACCTEST}');"
  echo -n 'terraform_token.save!;'
) |
docker exec -i gitlab gitlab-rails console


