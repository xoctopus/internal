package main

/*
#cgo CFLAGS: -I.
#include "shared.h"
*/
import "C"

import (
	"fmt"
	"structs"
	"unsafe"
)

// Packet contains POD(plain-old-data) fields only. Complex data type is not
// supported, like string, slice, map, interface and pointer. If memory lifetime
// is managed by golang garbage collection, DON'T USE THAT.
type Packet struct {
	_      structs.HostLayout // <-- keep host ABI layout
	ID     uint32
	Flags  uint16
	Length uint16
}

func main() {
	var p Packet

	C.fill_packet((*C.struct_Packet)(unsafe.Pointer(&p)), 1, 2, 3)

	fmt.Println("After C.fill_packet:")
	fmt.Printf("ID     = %d\n", p.ID)
	fmt.Printf("Flags  = %d\n", p.Flags)
	fmt.Printf("Length = %d\n", p.Length)
}
