TOOLCHAIN=/usr/local/android-ndk-r19c/toolchains/llvm/prebuilt/linux-x86_64
ARCH=armv7a
export AR=$TOOLCHAIN/bin/$ARCH-linux-android-ar
export AS=$TOOLCHAIN/bin/$ARCH-linux-android-as
export CC=$TOOLCHAIN/bin/$ARCH-linux-androideabi28-clang
export CXX=$TOOLCHAIN/bin/$ARCH-linux-androideabi28-clang++
export LD=$TOOLCHAIN/bin/$ARCH-linux-android-ld
export RANLIB=$TOOLCHAIN/bin/$ARCH-linux-android-ranlib
export STRIP=$TOOLCHAIN/bin/$ARCH-linux-android-strip

#GOOS=android GOARCH=arm CGO_ENABLED=1 gobind -lang=go,java -outdir=./build/bind github.com/mysteriumnetwork/go-openvpn/test
#GOOS=android GOARCH=arm CGO_ENABLED=1 go build -buildmode=c-shared -o lib.o ./test/test.go
GOOS=android GOARCH=arm CGO_ENABLED=1 gomobile bind -target=android/arm github.com/mysteriumnetwork/go-openvpn/test