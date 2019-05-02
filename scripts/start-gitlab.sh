#!/bin/bash -e
if [[ $MAKE_TARGET != "testacc" ]]; then echo "not starting gitlab!"; exit 0; fi
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
docker exec -i gitlab gitlab-rails console < scripts/generate-access-token.rb
