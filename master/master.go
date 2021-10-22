package master

import (
	"fmt"
	"log"
	"net/http"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tzY15368/lazarus/config"
	"github.com/tzY15368/lazarus/gen-go/RPCService"
	"github.com/tzY15368/lazarus/master/handlers"
	"github.com/tzY15368/lazarus/master/handlers/rpc"
	"github.com/tzY15368/lazarus/master/models"
)

// gin server for handling business
var externalG *gin.Engine

func say(ctx *gin.Context) {
	ctx.String(http.StatusOK, "helo")
}

func StartMaster() {
	cfg := config.Cfg.Master
	// seteup db
	err := models.SetupDB(cfg.Db)
	if err != nil {
		log.Fatal("db conn error", err)
	}

	// init web server
	webAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	externalG = gin.Default()
	gin.SetMode(gin.ReleaseMode)
	sessionStore := cookie.NewStore([]byte(cfg.Secret))
	externalG.Use(sessions.Sessions("masterSession", sessionStore))
	lazarus := externalG.Group("/lazarus")
	lazarus.StaticFile("/", "html/index.html")
	lazarus.StaticFS("/home", http.Dir("html/panel/build"))
	lazarus.GET("/login", handlers.LoginHandler)
	lazarus.GET("/logout", handlers.LogoutHandler)
	lazarus.GET("/update", handlers.LoginRequired, handlers.UpdateSubscription)
	lazarus.GET("/user", handlers.LoginRequired, handlers.UserInfoHandler)
	lazarus.GET("/s/json", handlers.LoginRequired, handlers.HandleSubscriptionJSON)
	lazarus.GET("/s/:token", handlers.HandleSubscription)
	lazarus.GET("/flush", handlers.LoginRequired, handlers.HandleTokenRefresh)
	go externalG.Run(webAddr)
	logrus.Info("started master web server at", webAddr)

	// init rpc server
	rpcAddr := fmt.Sprintf("%s:%d", cfg.RpcHost, cfg.RpcPort)
	transport, err := thrift.NewTServerSocket(rpcAddr)
	if err != nil {
		logrus.Fatal(err)
	}
	processor := RPCService.NewLazarusServiceProcessor(&rpc.LazarusService{})
	transportFactory := thrift.NewTTransportFactory()
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)
	go server.Serve()
	logrus.Info("started master RPC server at", rpcAddr)
}
