package protocol

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/tuanliang/restful-api-demo/apps"
	"github.com/tuanliang/restful-api-demo/conf"
)

func NewHttpService() *HttpService {
	r := gin.Default()
	server := &http.Server{
		ReadHeaderTimeout: 60 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1M
		Addr:              conf.C().App.HttpAddr(),
		Handler:           r,
	}
	return &HttpService{
		server: server,
		l:      zap.L().Named("HTTP Service"),
		r:      r,
	}
}

type HttpService struct {
	server *http.Server
	l      logger.Logger
	r      gin.IRouter
}

func (s *HttpService) Start() error {
	// 加载handler,把所有模块的Handler注册给了Gin Router
	apps.InitGin(s.r)

	// 已加载App的日志信息
	apps := apps.LoadedGinApps()
	s.l.Info("loaded gin apps :%v", apps)

	if err := s.server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed { // 正常退出
			s.l.Info("service stopped success")
			return nil
		}
		return fmt.Errorf("start service error,%s", err.Error())
	}

	return nil
}
func (s *HttpService) Stop() {
	s.l.Info("start graceful shundown")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		s.l.Warnf("shut down http service error,%s", err)
	}
}
