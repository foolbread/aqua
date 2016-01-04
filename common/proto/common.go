//@auther foolbread
//@time 2015-11-13
//@file aqua/common/proto/common.go
package proto

import (
	"encoding/binary"
	"time"
)

//status
const (
	STATUS_OK     = 0
	ALREADY_LOGIN = 1
	SERVICE_ERROR = 2
)

//cmd
const (
	KEEPALIVE_CMD      = 0
	LOGINREQ_CMD       = 1
	LOGINRES_CMD       = 2
	REDIRECT_CMD       = 4
	SERVICEREQ_CMD     = 5
	SERVICERES_CMD     = 6
	CONREGISTERREQ_CMD = 7
	CONREGISTERRES_CMD = 8
	LOCREGISTERREQ_CMD = 9
	LOCREGISTERRES_CMD = 10
)

var KeepAlive []byte = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0c, 0x00, 0x00, 0x00, 0x00}

func FillHead(d []byte, cmd uint32) []byte {
	var buf []byte = make([]byte, HEAD_LEN, 1024)

	binary.BigEndian.PutUint32(buf[4:], uint32(len(d))+HEAD_LEN)
	binary.BigEndian.PutUint32(buf[8:], cmd)

	buf = append(buf, d...)

	return buf
}

func UnmarshalLoginReq(d []byte) (*LoginRequest, error) {
	var req LoginRequest
	err := req.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func MarshalLoginRes(token []byte, status int32, cid string) ([]byte, error) {
	var res LoginResponse
	res.Cid = cid
	res.Token = token
	res.Status = status
	res.ServerTime = time.Now().Unix()

	d, err := res.Marshal()
	if err != nil {
		return nil, err
	}

	buf := FillHead(d, LOGINRES_CMD)

	return buf, nil
}

func MarshalRedirect(status int, token []byte, addr string) ([]byte, error) {
	var res RedirectResponse
	res.Status = int32(status)
	res.Token = token
	res.Addr = addr

	d, err := res.Marshal()
	if err != nil {
		return nil, err
	}

	buf := FillHead(d, REDIRECT_CMD)

	return buf, nil
}

func UnmarshalServiceReq(d []byte) (*ServiceRequest, error) {
	var req ServiceRequest

	err := req.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func UnmarshalServiceRes(d []byte) (*ServiceResponse, error) {
	var res ServiceResponse

	err := res.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func MarshalServiceRes(to []byte, t int32, sn string, s int32, data []byte) ([]byte, error) {
	var res ServiceResponse
	res.Token = to
	res.ServiceType = t
	res.Sn = sn
	res.Status = s
	res.Payload = data

	d, err := res.Marshal()
	if err != nil {
		return nil, err
	}

	buf := FillHead(d, SERVICERES_CMD)

	return buf, nil
}

func UnmarshalConnectRegisterReq(d []byte) (*ConnectRegisterReq, error) {
	var req ConnectRegisterReq
	err := req.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func MarshalConnectRegisterReq(id uint32, addr string) ([]byte, error) {
	var req ConnectRegisterReq
	req.Id = id
	req.ListenAddr = addr

	d, err := req.Marshal()
	if err != nil {
		return nil, err
	}

	buf := FillHead(d, CONREGISTERREQ_CMD)

	return buf, nil
}

func UnmarshalConnectRegisterRes(d []byte) (*ConnectRegisterRes, error) {
	var res ConnectRegisterRes

	err := res.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func MarshalConnectRegisterRes(s int32) ([]byte, error) {
	var res ConnectRegisterRes
	res.Status = s

	d, err := res.Marshal()
	if err != nil {
		return nil, err
	}

	buf := FillHead(d, CONREGISTERRES_CMD)

	return buf, nil
}

func UnmarshalLogicRegisterReq(d []byte) (*LogicRegisterReq, error) {
	var req LogicRegisterReq
	err := req.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func MarshalLogicRegisterReq(t uint32) ([]byte, error) {
	var req LogicRegisterReq
	req.ServiceType = t

	d, err := req.Marshal()
	if err != nil {
		return nil, err
	}

	buf := FillHead(d, LOCREGISTERREQ_CMD)

	return buf, nil
}

func UnmarshalLogicRegisterRes(d []byte) (*LogicRegisterRes, error) {
	var res LogicRegisterRes
	err := res.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func MarshalLogicRegisterRes(id uint32, s int32) ([]byte, error) {
	var res LogicRegisterRes
	res.Status = s
	res.Id = id

	d, err := res.Marshal()
	if err != nil {
		return nil, err
	}

	buf := FillHead(d, LOCREGISTERRES_CMD)

	return buf, nil
}
