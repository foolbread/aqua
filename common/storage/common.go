//@auther foolbread
//@time 2016-01-21
//@file aqua/common/storage/common.go
package storage

import (
	"crypto/md5"
)

func Md5ToByte(s string) byte {
	k := md5.Sum([]byte(s))
	return k[0]
}
