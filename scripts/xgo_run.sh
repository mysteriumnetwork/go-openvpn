#!/usr/bin/env bash
set -e

COMMAND=$@
if [ -z "$COMMAND" ]; then
    printf "\e[0;31m%s\e[0m\n" "Missing command!"
    exit 1
fi

docker run -it --rm \
    -v $PWD:/go-src-root/ \
    -v $PWD:/go/src/github.com/mysteriumnetwork/go-openvpn \
    -w /go/src/github.com/mysteriumnetwork/go-openvpn \
    --entrypoint "/bin/bash" \
    mysteriumnetwork/xgo:1.11 -c "${COMMAND}"
