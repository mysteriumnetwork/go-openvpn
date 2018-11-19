#!/usr/bin/env bash
set -e
. /go-src-root/scripts/helpers.sh

cd /go-src-root/build
rm -rf core
fetch_openvpn3

rm -rf /go-src-root/openvpn3/bridge/*.a /go-src-root/openvpn3/bridge/*.h
cp -f core/adapter/*.h /go-src-root/openvpn3/bridge/.

export O3=$PWD
export DEP_DIR=$PWD/dep_dir
mkdir -p $DEP_DIR
export DL=$PWD/dls
mkdir -p $DL

echo "Deps are in: $DEP_DIR"
echo "DLs go to: $DL"
echo "O3 is: $O3"

/go-src-root/scripts/x-compile-linux.sh
/go-src-root/scripts/x-compile-mac.sh
/go-src-root/scripts/x-compile-windows.sh
/go-src-root/scripts/x-compile-ios.sh
/go-src-root/scripts/x-compile-android-arm64.sh
/go-src-root/scripts/x-compile-android-armeabi-v7a.sh
/go-src-root/scripts/x-compile-android-x86.sh
/go-src-root/scripts/x-compile-android-amd64.sh


