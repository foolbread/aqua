//@auther foolbread
//@time 2016-01-21
//@file aqua/common/storage/common.go
package storage

import (
	aerr "aqua/common/error"
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/foolbread/fbcommon/golog"
)

func Md5ToByte(s string) byte {
	k := md5.Sum([]byte(s))
	return k[0]
}

func ParseSession(s string) (token []byte, id int) {
	idx := strings.LastIndex(s, "_")
	if idx < 0 {
		golog.Error(aerr.ErrSessionFormat)
		return
	}

	token, _ = hex.DecodeString(s[:idx])
	id, _ = strconv.Atoi(s[idx+1:])

	return token, id
}
