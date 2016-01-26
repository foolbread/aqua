//@auther foolbread
//@time 2015-11-13
//@file aqua/common/proto/common.go
package proto

import (
	"encoding/binary"
)

//status
const (
	STATUS_OK         = 0
	ALREADY_LOGIN     = 1
	SERVICE_ERROR     = 2
	MESSAGE_FULL      = 3
	WAITTING_RESPONSE = 4
)

//cmd
const (
	KEEPALIVE_CMD      = 0
	LOGINREQ_CMD       = 1
	LOGINRES_CMD       = 2
	REDIRECT_CMD       = 4
	CONREGISTERREQ_CMD = 5
	CONREGISTERRES_CMD = 6
	LOCREGISTERREQ_CMD = 7
	LOCREGISTERRES_CMD = 8

	LOGICSERVICEREQ_CMD = 9
	LOGICSERVICERES_CMD = 10
)

var KeepAlive []byte = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0c, 0x00, 0x00, 0x00, 0x00}

func FillHead(d []byte, cmd uint32) []byte {
	var buf []byte = make([]byte, HEAD_LEN, 1024)

	binary.BigEndian.PutUint32(buf[4:], uint32(len(d))+HEAD_LEN)
	binary.BigEndian.PutUint32(buf[8:], cmd)

	buf = append(buf, d...)

	return buf
}
