package router

import (
	"mayfly-go/internal/sys/api"
	"mayfly-go/internal/sys/application"
	"mayfly-go/pkg/req"

	"github.com/gin-gonic/gin"
)

func InitSyslogRouter(router *gin.RouterGroup) {
	s := &api.Syslog{
		SyslogApp: application.GetSyslogApp(),
	}
	sys := router.Group("syslogs")
	{
		sys.GET("", func(c *gin.Context) {
			req.NewCtxWithGin(c).Handle(s.Syslogs)
		})
	}
}
