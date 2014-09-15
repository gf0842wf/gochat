package main

import (
	"fmt"
)

// s: ">HIB12sQ"
// s: "<HIB12sQ"

func CalcSize(s string) (size int) {
	s = s[1:]
	for _, v := range s {
		switch v {

		}
	}

	return
}
func Unpack(s string, data []byte) (msg []interface{}) {
	dian := s[0]
	fmt.Println(dian)
	s = s[1:]

	msg = make([]interface{}, 2)

	for _, v := range s {
		if v == 'I' {
			msg[0] = 2
		}
	}
	return
}
func main() {
	s := ">IHQ12sB"
	data := []byte{1, 2, 1, 2, 3}
	Unpack(s, data)
}
