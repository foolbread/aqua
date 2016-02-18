//@auther foolbread
//@time 2016-01-25
//@file aqua/relation_server/server/relation_helper.go
package server

import (
	aproto "aqua/common/proto"
	astorage "aqua/common/storage"
	"aqua/relation_server/config"
	"aqua/relation_server/storage"
	"encoding/base64"
	"encoding/hex"
	"strings"

	"github.com/foolbread/fbcommon/golog"
)

func SendServiceMsg(con *connectServer, cid string, r *aproto.ServiceRequest, rr []byte) {
	data, err := aproto.MarshalServiceRes(r.Token, int32(config.GetConfig().GetServiceType()), r.Sn, aproto.STATUS_OK, rr)
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

func SendServiceMsgEx(cid string, rr []byte, sn string) {
	hnl := storage.GetStorage().GetSessionHandler(cid)

	online, err := hnl.IsExistSession(cid)
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

	data, err := aproto.MarshalServiceRes(token, int32(config.GetConfig().GetServiceType()), sn, aproto.STATUS_OK, rr)
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

///////////////////////////////////////////////////////////////////////////////////////////////////////
func SendMsg(con *connectServer, cid string, r *aproto.ServiceRequest, rp *aproto.RelationPacket, b bool /*cache flag*/) {
	data, err := rp.Marshal()
	if err != nil {
		golog.Error(err)
		return
	}

	if b {
		hnl := storage.GetStorage().GetSessionHandler(cid)
		err = hnl.AddRelationMsg(cid, base64.StdEncoding.EncodeToString(data), int(rp.Id))
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
