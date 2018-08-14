package process
// #include <stdlib.h>
import "C"
import "unsafe"

type Args struct {
	pointers []*C.char
}

func (a * Args) AddAll(args... string) {
	for _ , arg := range args {
		a.pointers = append(a.pointers, C.CString(arg))
	}
}

func (a * Args) cPointer() (**C.char) {
	return (**C.char)(&a.pointers[0])
}

func (a * Args) cCount() C.int {
	return C.int(len(a.pointers))
}

func (a * Args) Free() {
	for _ , arg := range a.pointers {
		C.free(unsafe.Pointer(arg))
	}
}