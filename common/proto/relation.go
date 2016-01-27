//@auther foolbread
//@time 2015-11-13
//@file aqua/common/proto/relation.go
package proto

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
