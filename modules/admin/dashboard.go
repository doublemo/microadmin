package admin


import (
	"net/http"
	
	"github.com/doublemo/msadmin/config"
	"github.com/gin-gonic/gin"
)

// dashboard 仪表盘
func dashboard(ctx *gin.Context, r *config.Registry) {
	ctx.HTML(http.StatusOK, "dashboard.tmpl", gin.H{
		"title": "Dashboard",
	})
}