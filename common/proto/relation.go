//@auther foolbread
//@time 2015-11-13
//@file aqua/common/proto/relation.go
package proto

//relation packet type
const (
	ADDFRIENDREQ_TYPE = iota + 1
	ADDFRIENDRES_TYPE
	DELFRIENDREQ_TYPE
	DELFRIENDRES_TYPE
	ADDBLACKREQ_TYPE
	ADDBLACKRES_TYPE
	DELBLACKREQ_TYPE
	DELBLACKRES_TYPE
	GETRMSGREQ_TYPE
	GETRMSGRES_TYPE
	RECVRMSGRES_TYPE
)

func UnmarshalAddFriendRes(d []byte) (*AddFriendRes, error) {
	var res AddFriendRes

	err := res.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func UnmarshalAddFriendReq(d []byte) (*AddFriendReq, error) {
	var req AddFriendReq
	err := req.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func UnmarshalDelFriendRes(d []byte) (*DelFriendRes, error) {
	var res DelFriendRes
	err := res.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func UnmarshalDelFriendReq(d []byte) (*DelFriendReq, error) {
	var req DelFriendReq
	err := req.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func UnmarshalAddBlackReq(d []byte) (*AddBlackReq, error) {
	var req AddBlackReq
	err := req.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func UnmarshalAddBlackRes(d []byte) (*AddBlackRes, error) {
	var res AddBlackRes
	err := res.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func UnmarshalDelBlackReq(d []byte) (*DelBlackReq, error) {
	var req DelBlackReq
	err := req.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func UnmarshalDelBlackRes(d []byte) (*DelBlackRes, error) {
	var res DelBlackRes
	err := res.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func UnmarshalRelationPacket(d []byte) (*RelationPacket, error) {
	var pa RelationPacket

	err := pa.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &pa, nil
}

func UnmarshalGetRMsgReq(d []byte) (*GetRelationPacketReq, error) {
	var req GetRelationPacketReq

	err := req.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func UnmarshalGetRMsgRes(d []byte) (*GetRelationPacketRes, error) {
	var res GetRelationPacketRes

	err := res.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func UnmarshalRecvRMsg(d []byte) (*RecvRelationPacket, error) {
	var res RecvRelationPacket

	err := res.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func MarshalRecvRMsg(cid string, ids []int64) ([]byte, error) {
	var res RecvRelationPacket
	res.Cid = cid
	res.Id = ids

	return res.Marshal()
}

func MarshalGetRMsgRes(msgs []*RelationPacket) ([]byte, error) {
	var res GetRelationPacketRes
	res.Msgs = msgs

	return res.Marshal()
}

func MarshalGetRMsgReq(cid string) ([]byte, error) {
	var req GetRelationPacketReq
	req.Cid = cid

	return req.Marshal()
}

func MarshalRelationPacket(ty int32, id int64, d []byte) ([]byte, error) {
	var pa RelationPacket
	pa.PacketType = ty
	pa.Id = id
	pa.Data = d

	return pa.Marshal()
}

func MarshalRelationPacketEx(ty int32, id int64, d []byte) *RelationPacket {
	var pa RelationPacket
	pa.PacketType = ty
	pa.Id = id
	pa.Data = d

	return &pa
}

func MarshalDelBlackRes(from string, black string, status int32) ([]byte, error) {
	var res DelBlackRes
	res.From = from
	res.Black = black
	res.Status = status

	return res.Marshal()
}

func MarshalDelBlackReq(from string, black string) ([]byte, error) {
	var req DelBlackReq
	req.From = from
	req.Black = black

	return req.Marshal()
}

func MarshalAddBlackRes(from string, black string, status int32) ([]byte, error) {
	var res AddBlackRes
	res.From = from
	res.Black = black
	res.Status = status

	return res.Marshal()
}

func MarshalAddBlackReq(from string, black string) ([]byte, error) {
	var req AddBlackReq
	req.From = from
	req.Black = black

	return req.Marshal()
}

func MarshalDelFriendReq(from string, friend string) ([]byte, error) {
	var req DelFriendReq
	req.From = from
	req.Friend = friend

	return req.Marshal()
}

func MarshalDelFriendRes(from string, friend string, status int32) ([]byte, error) {
	var res DelFriendRes
	res.From = from
	res.Friend = friend
	res.Status = status

	return res.Marshal()
}

func MarshalAddFriendReq(from string, friend string, d []byte) ([]byte, error) {
	var req AddFriendReq
	req.From = from
	req.Friend = friend
	req.Data = d

	return req.Marshal()
}

func MarshalAddFriendRes(from string, friend string, s int32) ([]byte, error) {
	var res AddFriendRes
	res.From = from
	res.Friend = friend
	res.Status = s

	return res.Marshal()
}
