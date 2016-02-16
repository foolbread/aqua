//@auther foolbread
//@time 2016-01-04
//@file aqua/singlechat_server/server/singlechat_helper.go
package server

import (
	aproto "aqua/common/proto"
	astorage "aqua/common/storage"
	"aqua/singlechat_server/config"
	"aqua/singlechat_server/storage"
	"encoding/base64"
	"encoding/hex"
	"strings"

	"github.com/foolbread/fbcommon/golog"
)

func SendServiceMsg(con *connectServer, cid string, r *aproto.ServiceRequest, pp []byte) {
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

	golog.Info("send service res to [connect_server]:", con.addr, "[id]:", con.id, "[cid]:", cid, "[token]:", strings.ToUpper(hex.EncodeToString(r.Token)), "[data_len]:", len(data))
}

func SendServiceMsgEx(cid string, pp []byte, sn string) {
	hnl := storage.GetStorage().GetSessionHandler(cid)
	online, err := hnl.IsExistSession(cid)
	if err != nil {
		golog.Error(err)
		return
	}

	if !online {
		return
	}

	//get usr session
	session, err := hnl.GetUsrSession(cid)
	if err != nil {
		golog.Error(err)
		return
	}

	token, id := astorage.ParseSession(session)
	if token == nil {
		return
	}

	data, err := aproto.MarshalServiceRes(token, int32(config.GetConfig().GetServiceType()), sn, aproto.STATUS_OK, pp)
	if err != nil {
		golog.Error(err)
		return
	}

	csvr := g_conmanager.getConnectSvr(id)

	if csvr != nil {
		err = csvr.SendToCon(data)
		if err != nil {
			golog.Error(err)
			return
		}
		golog.Info("send service res to [connect_server]:", csvr.addr, "[id]:", csvr.id, "[cid]:", cid, "[token]:", strings.ToUpper(hex.EncodeToString(token)), "[data_len]:", len(data))
	}

}

func SendMsg(con *connectServer, cid string, r *aproto.ServiceRequest, pp *aproto.PeerPacket, b bool /*cache flag*/) {
	data, err := pp.Marshal()
	if err != nil {
		golog.Error(err)
		return
	}

	if b {
		hnl := storage.GetStorage().GetSingleHandler(cid)
		err = hnl.AddPeerMsg(cid, base64.StdEncoding.EncodeToString(data), int(pp.Id))
		if err != nil {
			golog.Error(err)
			return
		}
	}

	if con != nil {
		SendServiceMsg(con, cid, r, data)
	} else {
		SendServiceMsgEx(cid, data, r.Sn)
	}
}
