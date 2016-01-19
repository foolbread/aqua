//@auther foolbread
//@time 2016-01-04
//@file aqua/singlechat_server/server/singlechat_helper.go
package server

import (
	aerr "aqua/common/error"
	aproto "aqua/common/proto"
	astorage "aqua/common/storage"
	"aqua/singlechat_server/config"
	"encoding/hex"
	"fbcommon/golog"
	"strconv"
	"strings"
)

func SendServiceMsg(con *connectServer, r *aproto.ServiceRequest, pp []byte) {
	data, err := aproto.MarshalServiceRes(r.Token, int32(config.GetConfig().GetServiceType()), r.Sn, aproto.STATUS_OK, pp)
	if err != nil {
		golog.Error(err)
		return
	}

	err = con.SendToCon(data)
	if err != nil {
		golog.Error(err)
		return
	}

	golog.Info("send service res to connect:", con.addr, "id:", con.id, "token:", hex.EncodeToString(r.Token), "data len:", len(data))
}

func SendPushMsgToUsr(r *aproto.ServiceRequest, rq *aproto.SendPeerMessageReq, handler *astorage.StorageHandler) {
	//get usr session
	session, err := handler.GetUsrSession(rq.Msg.To)
	if err != nil {
		golog.Error(err)
		return
	}

	idx := strings.LastIndex(session, "_")
	if idx < 0 {
		golog.Error(aerr.ErrSessionFormat)
		return
	}

	token, _ := hex.DecodeString(session[:idx])
	id, _ := strconv.Atoi(session[idx+1:])

	//construct PushPeerMsgReq
	d, err := aproto.MarshalPushPMsgReq(rq.Msg)
	if err != nil {
		golog.Error(err)
		return
	}

	//construct PeerPacket
	p, err := aproto.MarshalPeerPacket(aproto.PUSHMSGREQ_TYPE, d)
	if err != nil {
		golog.Error(err)
		return
	}

	//construct ServiceRes
	data, err := aproto.MarshalServiceRes(token, int32(config.GetConfig().GetServiceType()), r.Sn, aproto.STATUS_OK, p)
	if err != nil {
		golog.Error(err)
		return
	}

	csvr := g_conmanager.getConnectSvr(id)
	if csvr != nil {
		err := csvr.SendToCon(data)
		if err != nil {
			golog.Error(err)
		}
	}
}
