package openvpn3

/*
#include <version.hpp>

*/
import "C"

func Version() string {
	cStr, _ := C.getVersion()
	return C.GoString(cStr)
}
