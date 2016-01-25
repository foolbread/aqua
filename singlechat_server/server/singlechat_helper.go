//@auther foolbread
//@time 2016-01-04
//@file aqua/singlechat_server/server/singlechat_helper.go
package server

import (
	aerr "aqua/common/error"
	aproto "aqua/common/proto"
	"aqua/singlechat_server/config"
	"aqua/singlechat_server/storage"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/foolbread/fbcommon/golog"
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

	golog.Info("send service res to [connect_server]:", con.addr, "[id]:", con.id, "[token]:", strings.ToUpper(hex.EncodeToString(r.Token)), "[data_len]:", len(data))
}

func SendPMsgToUsr(r *aproto.ServiceRequest, rq *aproto.SendPeerMessageReq) {
	hnl := storage.GetStorage().GetSessionHandler(rq.Msg.To)
	//get usr session
	session, err := hnl.GetUsrSession(rq.Msg.To)
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

	//construct SendPeerMsgReq
	d, err := rq.Marshal()
	if err != nil {
		golog.Error(err)
		return
	}

	//construct PeerPacket
	p, err := aproto.MarshalPeerPacket(aproto.SENDPMSGREQ_TYPE, d)
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
		golog.Info("Send Peer Msg to [cid]:", rq.Msg.To, "[token]:", session[:idx], "[connect_id]:", id, "[data_len]:", len(data))
	}
}
