//@auther foolbread
//@time 2015-11-13
//@file aqua/common/net/net.go
package net

import (
	aerror "aqua/common/error"
	aproto "aqua/common/proto"
	"net"
	"time"

	fbnet "github.com/foolbread/fbcommon/net"
)

func RecvPacket(c net.Conn, d []byte) (uint32, uint32, error) {
	err := fbnet.ReadByCnt(c, d[:aproto.HEAD_LEN])

	if err != nil {
		return 0, 0, err
	}

	h := aproto.UnmarshalHead(d[:aproto.HEAD_LEN])
	if h.Length > uint32(len(d)) {
		return 0, 0, aerror.ErrPacketLen
	}

	if h.Length > aproto.HEAD_LEN {
		err = fbnet.ReadByCnt(c, d[aproto.HEAD_LEN:h.Length])
		if err != nil {
			return 0, 0, err
		}
	}

	return h.Length, h.Cmd, nil
}

func RecvPacketEx(c net.Conn, d []byte, timeout time.Duration) (uint32, uint32, error) {
	err := fbnet.ReadByTimeout(c, d[:aproto.HEAD_LEN], timeout)
	if err != nil {
		return 0, 0, err
	}

	h := aproto.UnmarshalHead(d[:aproto.HEAD_LEN])
	if h.Length > uint32(len(d)) {
		return 0, 0, aerror.ErrPacketLen
	}

	if h.Length > aproto.HEAD_LEN {
		err = fbnet.ReadByTimeout(c, d[aproto.HEAD_LEN:h.Length], timeout)
		if err != nil {
			return 0, 0, err
		}
	}

	return h.Length, h.Cmd, nil
}

func SendPacket(c net.Conn, d []byte) error {
	return fbnet.WriteByCnt(c, d)
}
