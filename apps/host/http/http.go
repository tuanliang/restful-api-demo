package http

import (
	"github.com/gin-gonic/gin"
	"github.com/tuanliang/restful-api-demo/apps"
	"github.com/tuanliang/restful-api-demo/apps/host"
)

func NewHostHTTPHandler() *Handler {
	return &Handler{}
}

// 通过写一个实体类，把内部的接口通过hTTp协议暴露出去
// 需要依赖内部接口的实现
// 该实体类，会实现Gin的Http Handler
type Handler struct {
	svc host.Service
}

func (h *Handler) Config() {
	if apps.HostService == nil {
		panic("dependence host service is nil")
	}
	h.svc = apps.HostService
}
func (h *Handler) Registry(r gin.IRouter) {
	r.POST("/hosts", h.createHost)
}
