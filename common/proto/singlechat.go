//@auther foolbread
//@time 2016-01-22
//@file aqua/common/proto/singlechat.go
package proto

import (
	"time"
)

//peer packet type
const (
	SENDPMSGREQ_TYPE = iota + 1
	SENDPMSGRES_TYPE
	GETPMSGREQ_TYPE
	GETPMSGRES_TYPE
	RECVPMSGRES_TYPE
)

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
	var req SendPeerMessageReq
	err := req.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func UnmarshalSendPMsgRes(d []byte) (*SendPeerMessageRes, error) {
	var res SendPeerMessageRes
	err := res.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &res, nil
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

func MarshalGetPMsgRes(msgs []*PeerPacket) ([]byte, error) {
	var res GetPeerMessageRes
	res.Msgs = msgs

	return res.Marshal()
}

func MarshalPeerPacket(t int32, id int64, data []byte) ([]byte, error) {
	var pp PeerPacket
	pp.PacketType = t
	pp.Id = id
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

func MarshalSendPMsgRes(cid string, status int32, sn string) ([]byte, error) {
	var smsg SendPeerMessageRes
	smsg.Cid = cid
	smsg.Status = status
	smsg.Sn = sn

	return smsg.Marshal()
}
