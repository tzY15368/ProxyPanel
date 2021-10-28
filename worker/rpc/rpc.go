package rpc

import (
	"context"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/sirupsen/logrus"
	"github.com/tzY15368/lazarus/config"
	"github.com/tzY15368/lazarus/gen-go/RPCService"
	"github.com/tzY15368/lazarus/worker/auth"
	"github.com/tzY15368/lazarus/worker/initialize"
	"github.com/tzY15368/lazarus/worker/sysinfo"
)

var ctx = context.TODO()

var rpcClient *RPCService.LazarusServiceClient

var sessionID string

const Port = 443

func init() {
	transport, err := thrift.NewTSocket(config.Cfg.Worker.MasterAddr)
	if err != nil {
		logrus.Fatal("initialize transport error", err)
	}
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	if err := transport.Open(); err != nil {
		logrus.Fatal("initialize protocol error", err)
	}
	iprot := protocolFactory.GetProtocol(transport)
	oprot := protocolFactory.GetProtocol(transport)
	rpcClient = RPCService.NewLazarusServiceClient(thrift.NewTStandardClient(iprot, oprot))
}
func Startup() (e error) {
	e = mustInitializeServer()
	if e != nil {
		return
	}
	e = mustRegisterSelf()
	if e != nil {
		return
	}
	return
}
func mustInitializeServer() error {
	res, err := rpcClient.DoInitialize(ctx, &RPCService.InitializeRequest{
		Mac: sysinfo.GetMacAddr(),
	})
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"add":  res.Add,
		"host": res.Host,
		"port": Port,
	}).Info("got initialize params")

	err = initialize.InitializeComponents(res.Add, res.Host, Port)
	if err != nil {
		return err
	}
	return nil
}

func mustRegisterSelf() error {
	res, err := rpcClient.DoRegister(ctx, &RPCService.RegisterRequest{
		Mac: sysinfo.GetMacAddr(),
	})
	config.Cfg.HeartBeatErrorThres = int(res.HeartBeatErrorThres)
	config.Cfg.HeartBeatRateIntervalSec = int(res.HeartBeatRateIntervalSec)
	return err
}

func SendHeartBeat() error {
	cpu := sysinfo.GetCPUPercent()
	mem := sysinfo.GetMemPercent()
	res, err := rpcClient.DoHeartBeat(ctx, &RPCService.HeartbeatRequest{
		Mac: sysinfo.GetMacAddr(),
		CPU: &cpu,
		Mem: &mem,
	})
	if err != nil {
		return err
	}
	handleHeartbeatResponse(res)
	return nil
}

func handleHeartbeatResponse(res *RPCService.HeartbeatResponse) {
	if res.HasUpdate {
		logrus.Info("has update on config")
		r := make(map[string]struct{})
		for _, v := range res.Data {
			r[v] = struct{}{}
		}
		auth.SetMap(r)
	}
}
