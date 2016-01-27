//@auther foolbread
//@time 2016-01-22
//@file aqua/common/proto/server.go
package proto

import (
	"time"
)

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
