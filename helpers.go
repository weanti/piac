package main

import (
	"encoding/binary"
	"fmt"
)

/*****************
*    Convert	 *
******************/

func uint16ArrayToByteArray(numbers []uint16) []byte {
	var bytes []byte = make([]byte, len(numbers)*2)
	for i, v := range numbers {
		binary.LittleEndian.PutUint16(bytes[2*i:2*i+2], v)
	}
	return bytes
}

func byteArrayToUint16Array(bytes []byte) []uint16 {
	if len(bytes)%2 != 0 {
		fmt.Errorf("Array length should be even, but it is %d", len(bytes))
		return []uint16{}
	}
	var numbers []uint16 = make([]uint16, len(bytes)/2)
	for i := range bytes {
		if i%2 == 0 {
			numbers[i/2] = binary.LittleEndian.Uint16(bytes[i : i+2])
		}
	}
	return numbers
}
