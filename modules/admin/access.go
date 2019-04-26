package admin

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Access 访问控制
func Access(s *gin.Engine) gin.HandlerFunc {
	return  func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		uid := session.Get("UID")
		if m, ok := uid.(uint64); !ok || m == 0 {
			ctx.Request.URL.Path = "/login"
			s.HandleContext(ctx)
			return
		}
	}
}