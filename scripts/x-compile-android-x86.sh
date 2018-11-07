#!/usr/bin/env bash
set -e

PLATFORM=android
SCRIPTS=android
LIBSUFFIX=android_x86

echo "Building for: $PLATFORM"

echo "Bootstrapping Android NDK"

$ANDROID_NDK_ROOT/build/tools/make_standalone_toolchain.py --install-dir=/usr/$ANDROID_CHAIN_386 --api=16 --arch=x86


echo BUILD DEPS
pushd $DEP_DIR
rm -rf asio* lz4* mbedtls* #lzo* boost* minicrypto openssl* polarssl* snappy*
echo "******* ASIO"
$O3/core/deps/asio/build-asio
echo "******* MBEDTLS"
TARGETS=android-x86 $O3/core/scripts/android/build-mbedtls
echo "******* LZ4"
TARGETS=android-x86 $O3/core/scripts/android/build-lz4
popd

. core/vars/setpath
. core/vars/vars-android-x86

cd core/adapter
PROF=$PLATFORM MTLS=1 NOSSL=1 LZ4=1 ASIO=1 ECHO=1 CO=1 build library

LIB_OUT="/go-src-root/openvpn3/bridge/libopenvpn3_${LIBSUFFIX}.a"
rm -rf $LIB_OUT

libs=("$DEP_DIR/mbedtls/mbedtls-$PLATFORM/library/libmbedtls.a" \
"$DEP_DIR/lz4/lz4-$PLATFORM/lib/liblz4.a" )

MORE_LIBS=${libs[@]} /go-src-root/scripts/concat-libs.sh $LIB_OUT library.o

$RANLIB_CMD $LIB_OUT | grep -v "has no symbols" || true
