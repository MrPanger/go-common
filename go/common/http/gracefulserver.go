package http

import (
	"fmt"
	"go/common/log"
	"strconv"

	"net"
	"net/http"
	"time"

	"github.com/braintree/manners"
	"github.com/gin-gonic/gin"
)

//为业务开发提供HttpServer和Handler
type GracefulServer struct {
	port           int
	listener       net.Listener
	gracefulServer *manners.GracefulServer
	tcpProxyPort   int
	engine         http.Handler
}

//h best for it is *gin.Engine
func NewGracefulServer(port int, h http.Handler) *GracefulServer {
	proxy := new(GracefulServer)
	proxy.port = port
	proxy.engine = h
	return proxy
}

func (gs *GracefulServer) Start() error {
	addr := fmt.Sprintf("0.0.0.0:%d", gs.port)
	gs.gracefulServer = manners.NewWithServer(&http.Server{
		Addr:         addr,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
		Handler:      gs.engine,
	})

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	go gs.gracefulServer.Serve(listener)
	gs.listener = listener
	log.Warn("gracefulServer start.now listen on " + strconv.Itoa(gs.port))

	return nil
}

func (gs *GracefulServer) Close() {
	if gs.gracefulServer != nil {
		gs.gracefulServer.Close()
		gs.gracefulServer = nil
		gs.listener.Close()
	}
	log.Warn("gracefulServer stop success")
}

func NewGracefulServerHandler(middleWare ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	result := gin.New()
	result.Use(gin.Recovery())
	result.Use(middleWare...)
	return result
}
