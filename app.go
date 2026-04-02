package main

import (
	"cyblog/internal/domain"
	"cyblog/internal/route"
	"cyblog/pkg/infra"
	"cyblog/pkg/log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/pprof"
	"github.com/spf13/viper"
)
import "github.com/gin-gonic/gin"

type MainApp struct {
	HttpServer      *http.Server
	Engine          *gin.Engine
	RootRouterGroup *gin.RouterGroup
	ServiceHub      *domain.ServiceHub
	port            uint
	host            string
	data            *infra.Data

	RegisterFunc route.RegisterFunc
}

func NewMainApp(
	vc *viper.Viper,
	hub *domain.ServiceHub,
	registerFunc route.RegisterFunc,
	registeredMiddleWire route.RegisteredMiddleWire,
	data *infra.Data,
) *MainApp {
	gin.SetMode(gin.DebugMode)

	e := gin.New()
	e.Use(gin.LoggerWithWriter(log.GetLogWriter()))
	e.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		log.SugaredLogger().Error("发生Panic!")
		log.SugaredLogger().Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}))
	registeredMiddleWire.Register()
	pprof.Register(e)
	registerFunc(e, hub)

	//route.RegisterNatsListener(data.NatsMQ.GetConn(), hub)

	app := &MainApp{
		Engine:       e,
		port:         vc.GetUint("server.http.port"),
		host:         vc.GetString("server.http.host"),
		ServiceHub:   hub,
		RegisterFunc: registerFunc,
	}
	app.PrintRoutes()
	return app
}

func (a *MainApp) PrintRoutes() {
	routes := a.Engine.Routes()
	log.SugaredLogger().Infof("Total routes: %d", len(routes))
	for _, route := range routes {
		log.SugaredLogger().Infof("Route: %-6s %s", route.Method, route.Path)
	}
}

func (a *MainApp) StartServer() error {

	addr := a.host + ":" + strconv.FormatUint(uint64(a.port), 10)
	log.GetLogger().Info("启动服务 " + addr)

	err := a.Engine.Run(addr)
	return err
}

func (a *MainApp) Close() error {
	return nil
}
