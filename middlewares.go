package msadmin

import (
	"github.com/doublemo/msadmin/config"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

//  middlewareSession session支持
func middlewareSession(r *config.Registry) {
	var (
		c = r.Configuration
	)

	sessionStore := memstore.NewStore([]byte(c.SessionSecret))
	r.Gin.Use(sessions.Sessions(c.SessionName, sessionStore))
}

// middlewareCSRF csrf支持
func middlewareCSRF(r *config.Registry) {
	var (
		c = r.Configuration
	)

	r.Gin.Use(csrf.Middleware(csrf.Options{
		Secret: c.CSRFSecret,
		ErrorFunc: func(ctx *gin.Context) {
			ctx.String(400, "CSRF token mismatch")
			ctx.Abort()
		},
	}))
}