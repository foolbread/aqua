//@auther foolbread
//@time 2015-11-13
//@file aqua/common/net/net.go
package net

import (
	aerror "aqua/common/error"
	aproto "aqua/common/proto"
	fbnet "fbcommon/net"
	"net"
	"time"
)

func RecvPacket(c net.Conn, d []byte) (uint32, error) {
	err := fbnet.ReadByCnt(c, d[:aproto.HEAD_LEN])

	if err != nil {
		return 0, err
	}

	h := aproto.UnmarshalHead(d[:aproto.HEAD_LEN])
	if h.Length > uint32(len(d[aproto.HEAD_LEN:])) {
		return 0, aerror.ErrPacketLen
	}

	err = fbnet.ReadByCnt(c, d[aproto.HEAD_LEN:h.Length-aproto.HEAD_LEN])
	if err != nil {
		return 0, err
	}

	return h.Length, nil
}

func RecvPacketEx(c net.Conn, d []byte, timeout time.Duration) (uint32, error) {
	err := fbnet.ReadByTimeout(c, d[:aproto.HEAD_LEN], timeout)
	if err != nil {
		return 0, err
	}

	h := aproto.UnmarshalHead(d[:aproto.HEAD_LEN])
	if h.Length > uint32(len(d[aproto.HEAD_LEN:])) {
		return 0, aerror.ErrPacketLen
	}

	err = fbnet.ReadByTimeout(c, d[aproto.HEAD_LEN:h.Length-aproto.HEAD_LEN], timeout)
	if err != nil {
		return 0, err
	}

	return h.Length, nil
}

func SendPacket(c net.Conn, d []byte) error {
	return fbnet.WriteByCnt(c, d)
}
