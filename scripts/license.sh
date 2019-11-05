#!/bin/bash

_license_dir=${GITLAB_LICENSE_DIR:-.license}
_license_file=${GITLAB_LICENSE_FILE:-JulienPivotto.gitlab-license}

mkdir -p "${_license_dir}";

if test -n "$encrypted_4a4b21de8599_key"; then
  echo decrypt

  openssl aes-256-cbc \
    -K "$encrypted_4a4b21de8599_key" \
    -iv "$encrypted_4a4b21de8599_iv" \
    -in "${_license_file}.enc" \
    -out "${_license_dir}/${_license_file}" \
    -d

  chmod 666 "${_license_dir}/${_license_file}"
fi
