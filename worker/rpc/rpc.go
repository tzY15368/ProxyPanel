package rpc

import (
	"context"
	"fmt"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/sirupsen/logrus"
	"github.com/tzY15368/lazarus/config"
	"github.com/tzY15368/lazarus/gen-go/RPCService"
)

var ctx = context.TODO()

func getRegistryClient() (*RPCService.LazarusServiceClient, error) {
	addr := fmt.Sprintf("%s:%d", config.Cfg.Worker.MasterIP, config.Cfg.Worker.MasterPort)
	transport, err := thrift.NewTSocket(addr)
	if err != nil {
		return nil, err
	}
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	if err := transport.Open(); err != nil {
		return nil, err
	}
	iprot := protocolFactory.GetProtocol(transport)
	oprot := protocolFactory.GetProtocol(transport)
	client := RPCService.NewLazarusServiceClient(thrift.NewTStandardClient(iprot, oprot))
	return client, nil
}

func RegisterSelf() error {
	client, err := getRegistryClient()
	if err != nil {
		return err
	}
	res, err := client.DoRegisterServer(ctx, &RPCService.RegisterRequest{
		Add:  "123",
		Host: "456",
		Ps:   "bbb",
	})
	if err != nil {
		panic(err)
		return err
	}
	_ = res
	logrus.Info("res", res)
	return nil
}
