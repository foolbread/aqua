//@auther foolbread
//@time 2015-11-13
//@file aqua/common/proto/relation.go
package proto

func UnmarshalAddFriendReq(d []byte) (*AddFriendReq, error) {
	var req AddFriendReq
	err := req.Unmarshal(d)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func MarshalAddFriendRes() ([]byte, error) {
	return nil, nil
}
