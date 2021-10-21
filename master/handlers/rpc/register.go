package rpc

import (
	"context"

	"github.com/tzY15368/lazarus/gen-go/RPCService"
)

type LazarusService struct {
}

func (ls *LazarusService) DoRegisterServer(ctx context.Context, rr *RPCService.RegisterRequest) (_r *RPCService.HeartbeatResponse, _err error) {
	a := "123"
	_r.HasUpdate = true
	_r.SessionID = &a
	_r.Data = RPCService.UserData{"123", "345"}
	return
}

func (ls *LazarusService) DoHeartBeat(ctx context.Context, hbr *RPCService.HeartbeatRequest) (res *RPCService.HeartbeatResponse, err error) {
	return
}
