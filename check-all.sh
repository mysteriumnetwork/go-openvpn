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

./xgo-check.sh --targets=linux/amd64,darwin/amd64,windows/amd64 --out=build/desktop $GOPATH/src/github.com/mysteriumnetwork/go-openvpn/examples/desktop
./xgo-check.sh --targets=ios/*,android/* --out=build/mobile $GOPATH/src/github.com/mysteriumnetwork/go-openvpn/examples/mobile

verifyAndroidLib build/mobile.aar x86 x86_64 arm64-v8a armeabi-v7a