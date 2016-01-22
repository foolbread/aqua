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
	MESSAGE_FULL  = 3
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

//peer packet type
const (
	SENDMSGREQ_TYPE = iota + 1
	SENDMSGRES_TYPE
	GETMSGREQ_TYPE
	GETMSGRES_TYPE
	PUSHMSGREQ_TYPE
	RECVMSGRES_TYPE
)

//relation packet type
const (
	ADDFRIENDREQ_TYPE = iota + 1
	ADDFRIENDRES_TYPE
	DELFRIENDREQ_TYPE
	DELFRIEDNRES_TYPE
	ADDBLACKREQ_TYPE
	ADDBLACKRES_TYPE
	DELBLACKREQ_TYPE
	DELBLACKRES_TYPE
	GETSYSMSGREQ_TYPE
	GETSYSMSGRES_TYPE
	RECVSYSMSGRES_TYPE
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

func UnmarshalLoginRes(d []byte) (*LoginResponse, error) {
	var res LoginResponse

	err := res.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &res, nil
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

func UnmarshalConnectRegisterReq(d []byte) (*ConnectRegisterReq, error) {
	var req ConnectRegisterReq
	err := req.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func UnmarshalConnectRegisterRes(d []byte) (*ConnectRegisterRes, error) {
	var res ConnectRegisterRes

	err := res.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func UnmarshalLogicRegisterReq(d []byte) (*LogicRegisterReq, error) {
	var req LogicRegisterReq
	err := req.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func UnmarshalLogicRegisterRes(d []byte) (*LogicRegisterRes, error) {
	var res LogicRegisterRes
	err := res.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func UnmarshalRedirectRes(d []byte) (*RedirectResponse, error) {
	var res RedirectResponse

	err := res.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func MarshalLoginReq(cid string, dt int32, ver string, token []byte) ([]byte, error) {
	var req LoginRequest
	req.Cid = cid
	req.DeviceType = dt
	req.ClientVersion = ver
	req.Token = token

	d, err := req.Marshal()
	if err != nil {
		return nil, err
	}

	buf := FillHead(d, LOGINREQ_CMD)

	return buf, nil
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

func MarshalServiceReq(to []byte, t int32, sn string, data []byte) ([]byte, error) {
	var req ServiceRequest
	req.Token = to
	req.ServiceType = t
	req.Sn = sn
	req.Payload = data

	d, err := req.Marshal()
	if err != nil {
		return nil, err
	}

	buf := FillHead(d, LOGICSERVICEREQ_CMD)

	return buf, nil
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

	buf := FillHead(d, LOGICSERVICERES_CMD)

	return buf, nil
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

////////////////////////////////////////////////////////////////////////
func UnmarshalPeerPacket(d []byte) (*PeerPacket, error) {
	var pg PeerPacket
	err := pg.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &pg, nil
}

func UnmarshalPeerMessage(d []byte) (*PeerMessage, error) {
	var pmsg PeerMessage
	err := pmsg.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &pmsg, nil
}

func UnmarshalSendPMsgReq(d []byte) (*SendPeerMessageReq, error) {
	var smsg SendPeerMessageReq
	err := smsg.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &smsg, nil
}

func UnmarshalGetPMsgReq(d []byte) (*GetPeerMessageReq, error) {
	var req GetPeerMessageReq
	err := req.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func UnmarshalRecvMsgRes(d []byte) (*RecvPeerMessageRes, error) {
	var res RecvPeerMessageRes
	err := res.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func UnmarshalGetPMsgRes(d []byte) (*GetPeerMessageRes, error) {
	var res GetPeerMessageRes
	err := res.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func UnmarshalPushMsgReq(d []byte) (*PushPeerMessageReq, error) {
	var req PushPeerMessageReq
	err := req.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func MarshalRecvMsgRes(cid string, ids []int64) ([]byte, error) {
	var res RecvPeerMessageRes
	res.Cid = cid
	res.Id = ids

	return res.Marshal()
}

func MarshalGetPMsgReq(cid string) ([]byte, error) {
	var req GetPeerMessageReq
	req.Cid = cid

	return req.Marshal()
}

func MarshalGetPMsgRes(msgs []*PeerMessage) ([]byte, error) {
	var res GetPeerMessageRes
	res.Msgs = msgs

	return res.Marshal()
}

func MarshalPeerPacket(t int32, data []byte) ([]byte, error) {
	var pp PeerPacket
	pp.PacketType = t
	pp.Data = data

	return pp.Marshal()
}

func MarshalSendMsgReq(from string, to string, data []byte) ([]byte, error) {
	var req SendPeerMessageReq
	req.Msg = new(PeerMessage)
	req.Msg.From = from
	req.Msg.To = to
	req.Msg.Time = time.Now().Unix()
	req.Msg.Data = data

	return req.Marshal()
}

func MarshalSendPMsgRes(status int32, sn string) ([]byte, error) {
	var smsg SendPeerMessageRes
	smsg.Status = status
	smsg.Sn = sn

	return smsg.Marshal()
}

func MarshalPushPMsgReq(msg *PeerMessage) ([]byte, error) {
	var pmsg PushPeerMessageReq
	pmsg.Msg = msg

	return pmsg.Marshal()
}

////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////
