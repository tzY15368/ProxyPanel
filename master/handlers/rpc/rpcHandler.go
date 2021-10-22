package rpc

import (
	"context"

	"github.com/tzY15368/lazarus/gen-go/RPCService"
	"github.com/tzY15368/lazarus/master/auth"
	"github.com/tzY15368/lazarus/master/handlers/servers"
)

type LazarusService struct {
}

func (ls *LazarusService) DoRegisterServer(ctx context.Context, rr *RPCService.RegisterRequest) (_r *RPCService.HeartbeatResponse, _err error) {
	sessionId := servers.RegisterServer(rr.Add, rr.Host, rr.Ps)
	auth.ChangeAuthMap()
	_r = &RPCService.HeartbeatResponse{
		HasUpdate: true,
		SessionID: &sessionId,
		Data:      auth.GetUserMap(),
	}
	return
}

func (ls *LazarusService) DoHeartBeat(ctx context.Context, hbr *RPCService.HeartbeatRequest) (res *RPCService.HeartbeatResponse, err error) {

	sessionID := hbr.SessionID
	servers.RegisterHeartbeat(sessionID)
	authMapDidChange := auth.AuthMapDidChange()
	res = &RPCService.HeartbeatResponse{
		HasUpdate: authMapDidChange,
		SessionID: &sessionID,
	}
	if authMapDidChange {
		res.Data = auth.GetUserMap()
	} else {
		res.Data = RPCService.UserData{}
	}
	return
}
