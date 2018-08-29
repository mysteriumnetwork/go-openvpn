#!/usr/bin/env bash
set -e
#x-compile.sh <profile>
#where profile can be linux | osx | mingw (windows)

PLATFORM=ios
SCRIPTS=mac
LIBSUFFIX=ios_arm64

echo "Building for: $PLATFORM"

BUILD_IOS=1 "core/scripts/$SCRIPTS/build-all"

. core/vars/setpath


cd core/adapter
PROF=$PLATFORM MTLS=1 NOSSL=1 LZ4=1 ASIO=1 ECHO=1 CO=1 OBJC=1 build library

LIB_OUT="/go-src-root/openvpn3/bridge/libopenvpn3_${LIBSUFFIX}.a"
rm -rf $LIB_OUT

x86_64-apple-darwin15-libtool -static -o $LIB_OUT \
    $DEP_DIR/mbedtls/mbedtls-$PLATFORM/library/libmbedtls.a \
    $DEP_DIR/lz4/lz4-$PLATFORM/lib/liblz4.a \
    library.o
