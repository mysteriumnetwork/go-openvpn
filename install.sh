TOOLCHAIN=/usr/local/android-ndk-r19c/toolchains/llvm/prebuilt/linux-x86_64
ARCH=i686
export AR=$TOOLCHAIN/bin/$ARCH-linux-android-ar
export AS=$TOOLCHAIN/bin/$ARCH-linux-android-as
export CC=$TOOLCHAIN/bin/$ARCH-linux-android28-clang
export CXX=$TOOLCHAIN/bin/$ARCH-linux-android28-clang++
export LD=$TOOLCHAIN/bin/$ARCH-linux-android-ld
export RANLIB=$TOOLCHAIN/bin/$ARCH-linux-android-ranlib
export STRIP=$TOOLCHAIN/bin/$ARCH-linux-android-strip
GOOS=android GOARCH=386 CGO_ENABLED=1 go install std


ARCH=armv7a
export AR=$TOOLCHAIN/bin/$ARCH-linux-android-ar
export AS=$TOOLCHAIN/bin/$ARCH-linux-android-as
export CC=$TOOLCHAIN/bin/$ARCH-linux-androideabi28-clang
export CXX=$TOOLCHAIN/bin/$ARCH-linux-androideabi28-clang++
export LD=$TOOLCHAIN/bin/$ARCH-linux-android-ld
export RANLIB=$TOOLCHAIN/bin/$ARCH-linux-android-ranlib
export STRIP=$TOOLCHAIN/bin/$ARCH-linux-android-strip
GOOS=android GOARCH=arm CGO_ENABLED=1 go install std

ARCH=aarch64
export AR=$TOOLCHAIN/bin/$ARCH-linux-android-ar
export AS=$TOOLCHAIN/bin/$ARCH-linux-android-as
export CC=$TOOLCHAIN/bin/$ARCH-linux-android28-clang
export CXX=$TOOLCHAIN/bin/$ARCH-linux-android28-clang++
export LD=$TOOLCHAIN/bin/$ARCH-linux-android-ld
export RANLIB=$TOOLCHAIN/bin/$ARCH-linux-android-ranlib
export STRIP=$TOOLCHAIN/bin/$ARCH-linux-android-strip
GOOS=android GOARCH=arm64 CGO_ENABLED=1 go install std

# ARCH=i686
# export AR=$TOOLCHAIN/bin/$ARCH-linux-android-ar
# export AS=$TOOLCHAIN/bin/$ARCH-linux-android-as
# export CC=$TOOLCHAIN/bin/$ARCH-linux-android28-clang
# export CXX=$TOOLCHAIN/bin/$ARCH-linux-android28-clang++
# export LD=$TOOLCHAIN/bin/$ARCH-linux-android-ld
# export RANLIB=$TOOLCHAIN/bin/$ARCH-linux-android-ranlib
# export STRIP=$TOOLCHAIN/bin/$ARCH-linux-android-strip
# GOOS=android GOARCH=386 CGO_ENABLED=1 go install std
