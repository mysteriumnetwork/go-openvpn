#!/usr/bin/env bash
set -e

. /go-src-root/scripts/helpers.sh
rm -rf core
fetch_openvpn3

cp -f core/adapter/*.h /go-src-root/openvpn3/bridge/.

export O3=`pwd`
export DEP_DIR=`pwd`/dep_dir
mkdir -p $DEP_DIR
export DL=/build/dls
mkdir -p $DL

echo "Deps are in: $DEP_DIR"
echo "DLs go to: $DL"
echo "O3 is: $O3"

/go-src-root/scripts/x-compile-linux.sh
/go-src-root/scripts/x-compile-mac.sh
/go-src-root/scripts/x-compile-windows.sh
/go-src-root/scripts/x-compile-ios.sh
/go-src-root/scripts/x-compile-android.sh
