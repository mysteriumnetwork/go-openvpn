package openvpn3

/*
#include <version.hpp>

*/
import "C"

// Version really is just a c++ function to force c++ style linking when (cross) compiling package
//TODO maybe we can safely remove all code from .ccp file and rename it dummy - just to kick c++ style linking
func Version() string {
	cStr, _ := C.getVersion()
	return C.GoString(cStr)
}
