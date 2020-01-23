#!/bin/bash -e

test "$MAKE_TARGET" == "testacc" || { echo "not starting gitlab!"; exit 0; }
test -z "$(docker ps -f 'name=gitlab' -q)" || { echo "Gitlab already running"; exit 0; }

_gitlab_image="gitlab/gitlab-ce"
extra=""
if test -n "${GITLAB_LICENSE_FILE}"; then
  _license_dir=${GITLAB_LICENSE_DIR:-.license}
  _license_file=${GITLAB_LICENSE_FILE:-JulienPivotto.gitlab-license}

  if ! test -r "${_license_dir}/${_license_file}"
  then
      echo "No license at ${_license_dir}/${_license_file}"
      exit 1
  fi
  _gitlab_image=gitlab/gitlab-ee
  extra+="-v $PWD/${_license_dir}:/license"
  extra+=" -e GITLAB_LICENSE_FILE=/license/${_license_file}"
fi

echo "Starting gitlab container..."
(
  set -o xtrace
  docker run -d --rm \
    --name gitlab \
    -e "GITLAB_ROOT_PASSWORD=${GITLAB_ROOT_PASSWORD:-adminadmin}" \
    -p 127.0.0.1:8080:80 \
    $extra \
    "${GITLAB_IMAGE:-$_gitlab_image}"
)

echo -n "Waiting for gitlab to be ready "
i=1
until wget -t 1 127.0.0.1:8080 -O /dev/null -q
do
    sleep 1
    echo -n .
    if test $((i%3)) -eq 0; then echo -n ' '; fi
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
