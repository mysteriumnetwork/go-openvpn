#!/bin/bash

TOOLCHAIN=$ANDROID_NDK_HOME/toolchains/llvm/prebuilt/linux-x86_64

ARCH=aarch64
export CC=$TOOLCHAIN/bin/$ARCH-linux-android21-clang
export CXX=$TOOLCHAIN/bin/$ARCH-linux-android21-clang++
export CGO_CXXFLAGS="-v"
export CGO_CCFLAGS="-v"
# export AR=$TOOLCHAIN/bin/$ARCH-linux-android-ar
# export AS=$TOOLCHAIN/bin/$ARCH-linux-android-as
# export LD=$TOOLCHAIN/bin/$ARCH-linux-android-ld
# export RANLIB=$TOOLCHAIN/bin/$ARCH-linux-android-ranlib
# export STRIP=$TOOLCHAIN/bin/$ARCH-linux-android-strip
GOOS=android GOARCH=arm64 CGO_ENABLED=1 gomobile bind -target=android/arm64 github.com/mysteriumnetwork/go-openvpn/openvpn3
# GOOS=android GOARCH=arm64 CGO_ENABLED=1 gobind -outdir ./build/app -lang="go,java" github.com/mysteriumnetwork/go-openvpn/openvpn3

# ARCH=i686
# export CC=$TOOLCHAIN/bin/$ARCH-linux-android28-clang
# export CXX=$TOOLCHAIN/bin/$ARCH-linux-android28-clang++
# export CGO_CXXFLAGS="-v"
# export CGO_CCFLAGS="-v"
# GOOS=android GOARCH=386 CGO_ENABLED=1 gomobile bind -target=android/386 github.com/mysteriumnetwork/go-openvpn/openvpn3

# ARCH=armv7a
# export CC=$TOOLCHAIN/bin/$ARCH-linux-androideabi28-clang
# stat $CC
# export CXX=$TOOLCHAIN/bin/$ARCH-linux-androideabi28-clang++
# export AR=$TOOLCHAIN/bin/$ARCH-linux-android-ar
# export AS=$TOOLCHAIN/bin/$ARCH-linux-android-as
# export LD=$TOOLCHAIN/bin/$ARCH-linux-android-ld
# export RANLIB=$TOOLCHAIN/bin/$ARCH-linux-android-ranlib
# export STRIP=$TOOLCHAIN/bin/$ARCH-linux-android-strip
# GOOS=android GOARCH=arm CGO_ENABLED=1 gomobile bind -target=android/arm github.com/mysteriumnetwork/go-openvpn/openvpn3
