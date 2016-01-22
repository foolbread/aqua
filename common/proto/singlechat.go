//@auther foolbread
//@time 2016-01-22
//@file aqua/common/proto/singlechat.go
package proto

import (
	"time"
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
