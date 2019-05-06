package msadmin

import (
	"github.com/doublemo/msadmin/config"
	"github.com/doublemo/msadmin/modules/admin"
	"github.com/doublemo/msadmin/modules/captcha"
	"github.com/gin-gonic/gin"
)

// routes 系统保留路由
func routes(r *config.Registry) {
	r.Gin.GET("/", func(ctx *gin.Context) {
		ctx.Request.URL.Path = "/admin"
		r.Gin.HandleContext(ctx)
	})

	// module admin
	admin.Routes(r)

	captcha.Routes(r)
}

// loadRoutes 加载动态路由信息
func loadRoutes(r *config.Registry) error {
	return nil
}
