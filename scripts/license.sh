mkdir license
[[ -n $encrypted_4a4b21de8599_key ]] && echo decrypt
[[ -n $encrypted_4a4b21de8599_key ]] && openssl aes-256-cbc -K $encrypted_4a4b21de8599_key -iv $encrypted_4a4b21de8599_iv -in JulienPivotto.gitlab-license.enc -out license/JulienPivotto.gitlab-license -d
chmod 666 license/JulienPivotto.gitlab-license || true
