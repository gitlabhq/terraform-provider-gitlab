#!/bin/bash -e
echo "Starting gitlab container..."
docker run -d -e GITLAB_ROOT_PASSWORD=adminadmin --rm -p 127.0.0.1:8080:80 --name gitlab gitlab/gitlab-ce
echo -n "Waiting for gitlab to be ready "
i=1
until wget -t 1 127.0.0.1:8080 -O /dev/null -q 
do
    sleep 1
    echo -n .
    if [[ $((i%3)) == 0 ]]; then echo -n ' '; fi
    let i++
done
echo
echo "Creating access token"
echo 'PersonalAccessToken.create(user_id: 1, scopes: [:api, :read_user], name: :terraform, token: :ACCTEST).save!' | docker exec -i gitlab gitlab-rails console
