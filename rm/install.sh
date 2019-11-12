
TOOLCHAIN=$ANDROID_NDK_HOME/toolchains/llvm/prebuilt/linux-x86_64
echo "Bootstrapping android/armv7a..."
ARCH=armv7a
CC=$TOOLCHAIN/bin/$ARCH-linux-androideabi16-clang CXX=$TOOLCHAIN/bin/$ARCH-linux-androideabi16-clang++ GOOS=android GOARCH=arm CGO_ENABLED=1 go install std

# echo "Bootstrapping android/arm64..."
# ARCH=aarch64
# CC=$TOOLCHAIN/bin/$ARCH-linux-android21-clang CXX=$TOOLCHAIN/bin/$ARCH-linux-android21-clang++ GOOS=android GOARCH=arm64 CGO_ENABLED=1 go install std

# echo "Bootstrapping android/i686..."
# ARCH=i686
# CC=$TOOLCHAIN/bin/$ARCH-linux-android21-clang CXX=$TOOLCHAIN/bin/$ARCH-linux-android21-clang++ GOOS=android GOARCH=386 CGO_ENABLED=1 go install std

# echo "Bootstrapping android/am64..."
# ARCH=x86_64
# CC=$TOOLCHAIN/bin/$ARCH-linux-android21-clang CXX=$TOOLCHAIN/bin/$ARCH-linux-android21-clang++ GOOS=android GOARCH=amd64 CGO_ENABLED=1 go install std
