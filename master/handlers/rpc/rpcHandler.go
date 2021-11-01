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

func (ls *LazarusService) DoRegister(ctx context.Context, rr *RPCService.RegisterRequest) (*RPCService.RegisterResponse, error) {
	params, err := servers.RegisterServer(rr.IP)
	if err != nil {
		return nil, err
	}
	auth.ChangeAuthMap()
	_r := &RPCService.RegisterResponse{
		HeartBeatRateIntervalSec: int32(config.Cfg.HeartBeatRateIntervalSec),
		HeartBeatErrorThres:      int32(config.Cfg.HeartBeatErrorThres),
		Add:                      params.Add,
		Host:                     params.Host,
	}
	return _r, nil
}

func (ls *LazarusService) DoHeartBeat(ctx context.Context, hbr *RPCService.HeartbeatRequest) (*RPCService.HeartbeatResponse, error) {
	err := servers.RegisterHeartbeat(hbr.IP)
	if err != nil {
		return nil, err
	}
	authMapDidChange := auth.AuthMapDidChange()
	res := &RPCService.HeartbeatResponse{
		HasUpdate: authMapDidChange,
	}
	if authMapDidChange {
		res.Data = auth.GetUserMap()
	} else {
		res.Data = RPCService.UserData{}
	}
	return res, nil
}
