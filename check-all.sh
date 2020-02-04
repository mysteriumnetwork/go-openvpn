#!/usr/bin/env bash
set -e

function verifyAndroidLib {
    local file=$1
    shift;

    for arch in ${@}; do
        echo -n "Verifying $arch ... "
        found=`unzip -l $file | grep -e "\/$arch\/"` || true
        if [ "$found" = "" ]; then
            echo "FAIL"
            return 1
        fi
        echo "OK"
    done
    return 0
}

docker run --rm \
    -v "$PWD"/build:/build \
    -v "$GOPATH"/.xgo-cache:/deps-cache:ro \
    -v "$PWD":/source \
    -e OUT=desktop \
    -e FLAG_V=false \
    -e FLAG_X=false \
    -e FLAG_RACE=false \
    -e FLAG_BUILDMODE=default \
    -e TARGETS=linux/amd64,darwin/amd64,windows/amd64 \
    mysteriumnetwork/xgo:1.13.6 github.com/mysteriumnetwork/go-openvpn/examples/desktop


docker run --rm \
    -v "$PWD"/build:/build \
    -v "$GOPATH"/.xgo-cache:/deps-cache:ro \
    -v "$GOPATH"/src:/ext-go/1/src:ro \
    -e OUT=mobile \
    -e FLAG_V=false \
    -e FLAG_X=false \
    -e FLAG_RACE=false \
    -e FLAG_BUILDMODE=default \
    -e TARGETS=android/* \
    -e EXT_GOPATH=/ext-go/1 \
    -e GO111MODULE=off \
    mysteriumnetwork/xgomobile:1.13.6 github.com/mysteriumnetwork/go-openvpn/examples/mobile

verifyAndroidLib build/mobile.aar x86 x86_64 arm64-v8a armeabi-v7a
