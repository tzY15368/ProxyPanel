package rpc

import (
	"context"
	"fmt"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/sirupsen/logrus"
	"github.com/tzY15368/lazarus/config"
	"github.com/tzY15368/lazarus/gen-go/RPCService"
	"github.com/tzY15368/lazarus/worker/auth"
)

var ctx = context.TODO()

var rpcClient *RPCService.LazarusServiceClient

var sessionID string

func getRPCClient() error {
	addr := fmt.Sprintf("%s:%d", config.Cfg.Worker.MasterIP, config.Cfg.Worker.MasterPort)
	transport, err := thrift.NewTSocket(addr)
	if err != nil {
		return err
	}
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	if err := transport.Open(); err != nil {
		return err
	}
	iprot := protocolFactory.GetProtocol(transport)
	oprot := protocolFactory.GetProtocol(transport)
	rpcClient = RPCService.NewLazarusServiceClient(thrift.NewTStandardClient(iprot, oprot))
	return nil
}

func handleHeartbeatResponse(res *RPCService.HeartbeatResponse) {
	sessionID = *res.SessionID
	if res.HasUpdate {
		logrus.Info("has update on config")
		r := make(map[string]struct{})
		for _, v := range res.Data {
			r[v] = struct{}{}
		}
		auth.SetMap(r)
	}
}

func RegisterSelf() error {
	err := getRPCClient()
	if err != nil {
		return err
	}
	res, err := rpcClient.DoRegisterServer(ctx, &RPCService.RegisterRequest{
		Add:  config.Cfg.Worker.Address,
		Host: config.Cfg.Worker.Host,
		Ps:   config.Cfg.Worker.PS,
	})
	if err != nil {
		return err
	}
	logrus.WithField("sessionID", *res.SessionID).Info("got session id")
	handleHeartbeatResponse(res)

	return nil
}

func SendHeartBeat() error {
	res, err := rpcClient.DoHeartBeat(ctx, &RPCService.HeartbeatRequest{
		SessionID: sessionID,
	})
	if err != nil {
		return err
	}
	handleHeartbeatResponse(res)
	return nil
}
