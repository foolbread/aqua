//@auther foolbread
//@time 2015-11-13
//@file aqua/common/proto/head.go
package proto

import (
	"encoding/binary"
)

const HEAD_LEN = 12

type Head struct {
	Magic  uint32
	Length uint32
	Cmd    uint32
}

func UnmarshalHead(d []byte) *Head {
	h := new(Head)
	h.Magic = binary.BigEndian.Uint32(d[0:4])
	h.Length = binary.BigEndian.Uint32(d[4:8])
	h.Cmd = binary.BigEndian.Uint32(d[8:12])

	return h
}
