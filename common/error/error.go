//@auther foolbread
//@time 2015-11-13
//@file aqua/common/error/error.go
package error

import (
	"errors"
)

var (
	ErrPacketLen = errors.New("packet len is too long!")
	ErrConnExsit = errors.New("connect is exsit!")
)
