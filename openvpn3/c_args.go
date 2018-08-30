package openvpn3

// #include <stdlib.h>
import "C"
import "unsafe"

type CCharPointer struct {
	Ptr *C.char
}

func NewCharPointer(str string) *CCharPointer {
	return &CCharPointer{
		Ptr: C.CString(str),
	}
}

func (ptr *CCharPointer) Delete() {
	C.free(unsafe.Pointer(ptr.Ptr))
}

type CCharPointerArray struct {
	pointers []*C.char
}

func (a *CCharPointerArray) AddAll(args ...string) {
	for _, arg := range args {
		a.pointers = append(a.pointers, C.CString(arg))
	}
}

func (a *CCharPointerArray) cPointer() **C.char {
	return (**C.char)(&a.pointers[0])
}

func (a *CCharPointerArray) cCount() C.int {
	return C.int(len(a.pointers))
}

func (a *CCharPointerArray) Free() {
	for _, arg := range a.pointers {
		C.free(unsafe.Pointer(arg))
	}
}
