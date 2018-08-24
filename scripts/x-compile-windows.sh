#!/usr/bin/env bash
set -e
#x-compile.sh <profile>
#where profile can be linux | osx | mingw (windows)

PLATFORM=mingw
SCRIPTS=win
LIBSUFFIX=windows_amd64

echo "Building for: $PLATFORM"

"core/scripts/$SCRIPTS/build-all"

. core/vars/setpath
. core/vars/vars-$PLATFORM
cd core/adapter
PROF=$PLATFORM MTLS=1 NOSSL=1 LZ4=1 ASIO=1 ECHO=1 CO=1 build library

LIB_OUT="/go-src-root/openvpn3/bridge/libopenvpn3_${LIBSUFFIX}.a"
rm -rf $LIB_OUT

libs=("$DEP_DIR/mbedtls/mbedtls-$PLATFORM/library/libmbedtls.a" \
"$DEP_DIR/lz4/lz4-$PLATFORM/lib/liblz4.a" )

MORE_LIBS=${libs[@]} /go-src-root/scripts/concat-libs.sh $LIB_OUT library.o

$RANLIB_CMD $LIB_OUT | grep -v "has no symbols" || true
