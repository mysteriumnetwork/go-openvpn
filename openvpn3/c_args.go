package openvpn3

// #include <stdlib.h>
import "C"
import "unsafe"

type cCharPointer struct {
	Ptr *C.char
}

func newCharPointer(str string) *cCharPointer {
	return &cCharPointer{
		Ptr: C.CString(str),
	}
}

func (ptr *cCharPointer) delete() {
	C.free(unsafe.Pointer(ptr.Ptr))
}

type cCharPointerArray struct {
	pointers []*C.char
}

func (a *cCharPointerArray) addAll(args ...string) {
	for _, arg := range args {
		a.pointers = append(a.pointers, C.CString(arg))
	}
}

func (a *cCharPointerArray) cPointer() **C.char {
	return (**C.char)(&a.pointers[0])
}

func (a *cCharPointerArray) cCount() C.int {
	return C.int(len(a.pointers))
}

func (a *cCharPointerArray) free() {
	for _, arg := range a.pointers {
		C.free(unsafe.Pointer(arg))
	}
}
