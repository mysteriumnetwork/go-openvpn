package test

// #include <stdio.h>
import "C"

type cCharPointer struct {
	Ptr *C.char
}

func Hello() string {
	return "hi"
}

// func main() {

// }
