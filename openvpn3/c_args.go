/*
 * Copyright (C) 2018 The "MysteriumNetwork/go-openvpn" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
