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
	GETRELATIONMSGREQ_TYPE
	GETRELATIONMSGRES_TYPE
	RECVRELATIONMSGRES_TYPE
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

func UnmarshalRelationpacket(d []byte) (*RelationPacket, error) {
	var pa RelationPacket

	err := pa.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &pa, nil
}

func MarshalRelationPacket(ty int32, d []byte) ([]byte, error) {
	var pa RelationPacket
	pa.PacketType = ty
	pa.Data = d

	return pa.Marshal()
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

func MarshalAddFriendReq(id int64, from string, friend string, d []byte) ([]byte, error) {
	var req AddFriendReq
	req.Id = id
	req.From = from
	req.Friend = friend
	req.Data = d

	return req.Marshal()
}

func MarshalAddFriendRes(id int64, from string, friend string, s int32) ([]byte, error) {
	var res AddFriendRes
	res.Id = id
	res.From = from
	res.Friend = friend
	res.Status = s

	return res.Marshal()
}
