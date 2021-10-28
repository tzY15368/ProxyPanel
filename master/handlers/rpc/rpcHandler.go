package rpc

import (
	"context"

	"github.com/tzY15368/lazarus/config"
	"github.com/tzY15368/lazarus/gen-go/RPCService"
	"github.com/tzY15368/lazarus/master/auth"
	"github.com/tzY15368/lazarus/master/handlers/servers"
)

type LazarusService struct {
}

func DoInitialize(ctx context.Context, rr *RPCService.InitializeRequest) (res *RPCService.InitializeResponse, _err error) {
	return
}

func (ls *LazarusService) DoRegisterServer(ctx context.Context, rr *RPCService.RegisterRequest) (_r *RPCService.RegisterResponse, _err error) {
	_err = servers.RegisterServer(rr.IP)
	if _err != nil {
		return
	}
	auth.ChangeAuthMap()
	_r = &RPCService.RegisterResponse{
		HeartBeatRateIntervalSec: int32(config.Cfg.HeartBeatRateIntervalSec),
		HeartBeatErrorThres:      int32(config.Cfg.HeartBeatErrorThres),
	}
	return
}

func (ls *LazarusService) DoHeartBeat(ctx context.Context, hbr *RPCService.HeartbeatRequest) (res *RPCService.HeartbeatResponse, err error) {
	servers.RegisterHeartbeat(hbr.IP)
	authMapDidChange := auth.AuthMapDidChange()
	res = &RPCService.HeartbeatResponse{
		HasUpdate: authMapDidChange,
	}
	if authMapDidChange {
		res.Data = auth.GetUserMap()
	} else {
		res.Data = RPCService.UserData{}
	}
	return
}
