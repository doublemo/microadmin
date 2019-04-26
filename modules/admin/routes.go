package admin

import (
	"github.com/doublemo/msadmin/config"
	"github.com/gin-gonic/gin"
)

func Rotues(r *config.Registry) {
	r.Gin.GET("/login", func(ctx *gin.Context) {
		login(ctx, r)
	})

	r.Gin.POST("/login", func(ctx *gin.Context) {
		postLogin(ctx, r)
	})

	accessRoute := r.Gin.Group("/admin", Access(r.Gin))
	accessRoute.GET("/", func(ctx *gin.Context) {
		dashboard(ctx, r)
	})

	accessRoute.GET("/logout", func(ctx *gin.Context) {
		logout(ctx, r)
	})
}