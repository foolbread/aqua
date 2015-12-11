//@auther foolbread
//@time 2015-11-13
//@file aqua/common/proto/common.go
package proto

import (
	"encoding/binary"
	gproto "github.com/golang/protobuf/proto"
	"time"
)

const LOGINREQ_CMD = 0x01
const LOGINRES_CMD = 0x02
const REDIRECT_CMD = 0x04
const SERVICEREQ_CMD = 0x05
const SERVICERES_CMD = 0x06

func UnmarshalLoginReq(d []byte) (*LoginRequest, error) {
	var req LoginRequest
	err := gproto.Unmarshal(d, &req)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func MarshalLoginRes(token []byte, status int32, cid string) ([]byte, error) {
	var buf []byte = make([]byte, HEAD_LEN, 1024)
	var res LoginResponse
	res.Cid = cid
	res.Token = token
	res.Status = status
	res.ServerTime = time.Now().Unix()

	d, err := gproto.Marshal(&res)
	if err != nil {
		return nil, err
	}

	binary.BigEndian.PutUint32(buf[4:], uint32(len(d)+HEAD_LEN))
	binary.BigEndian.PutUint32(buf[8:], LOGINRES_CMD)

	buf = append(buf, d...)
	return buf, nil
}

func MarshalRedirect(token []byte, addr string) ([]byte, error) {
	var buf []byte = make([]byte, HEAD_LEN, 1024)
	var res RedirectResponse
	res.Status = 0
	res.Token = token
	res.Addrs = addr
	d, err := gproto.Marshal(&res)
	if err != nil {
		return nil, err
	}

	binary.BigEndian.PutUint32(buf[4:], uint32(len(d)+HEAD_LEN))
	binary.BigEndian.PutUint32(buf[8:], REDIRECT_CMD)

	buf = append(buf, d...)
	return buf, nil
}

func UnmarshalServiceReq(d []byte) (*ServiceRequest, error) {
	var req ServiceRequest
	err := gproto.Unmarshal(d, &req)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func MarshalServiceRes(t int32, sn string, s int32, data []byte) ([]byte, error) {
	var buf []byte = make([]byte, HEAD_LEN, 1024)
	var res ServiceResponse
	res.ServiceType = t
	res.Sn = sn
	res.Status = s
	res.Payload = data

	d, err := gproto.Marshal(&res)
	if err != nil {
		return nil, err
	}

	binary.BigEndian.PutUint32(buf[4:], uint32(len(d)+HEAD_LEN))
	binary.BigEndian.PutUint32(buf[8:], SERVICERES_CMD)

	buf = append(buf, d...)
	return buf, nil
}
