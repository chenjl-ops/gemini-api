package server

import (
	"context"
	"errors"
	"fmt"
	log "log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type server struct {
	App *gin.Engine
	//Validate *validator.Validate
}

func NewServer() (*server, error) {
	//gin.ForceConsoleColor()
	return &server{
		App: gin.New(),
	}, nil
}

func StartServer() error {
	server, err1 := NewServer()
	if err1 != nil {
		return err1
	}

	// 初始化header
	server.initRequestHeader()
	// 初始化trace id
	server.initTraceId()
	// 初始化日志
	server.initLog()

	// 初始化swagger
	//server.InitSwagger()

	// 初始化路由
	server.InitRouter(true)
	server.LoadUrl()

	//启动服务
	//err := server.Run()

	//优雅启动
	err := server.GracefulRun()
	if err != nil {
		return err
	}
	return nil
}

// Run 启动服务
func (s *server) Run() error {
	port := fmt.Sprintf(":%d", 8888)
	return s.App.Run(port)
}

// GracefulRun 优雅启动服务
func (s *server) GracefulRun() error {
	port := fmt.Sprintf(":%d", 8888)

	srv := &http.Server{
		Addr:    port,
		Handler: s.App.Handler(),
	}
	log.Info("Start Service Successful Listen Port", port)

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Listen Service Failed", "err", err)
			os.Exit(1)
		} else {
			log.Info("Start Service Error: ", err)
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server Shutdown: ", err)
		os.Exit(-1)
	}
	select {
	case <-ctx.Done():
		fmt.Println("timeout of 5 seconds.")
	}
	fmt.Println("Server exiting")
	return nil
}
