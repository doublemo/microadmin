package captcha

import (
	"bytes"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/dchest/captcha"
	"github.com/doublemo/msadmin/config"
	"github.com/gin-gonic/gin"
)

func Routes(r *config.Registry) {
	r.Gin.GET("/captcha/:name", func(ctx *gin.Context) {
		ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		ctx.Header("Pragma", "no-cache")
		ctx.Header("Expires", "0")

		file := ctx.Param("name")
		ext := path.Ext(file)
		id := file[:len(file)-len(ext)]
		if ext == "" || id == "" {
			ctx.String(http.StatusNotFound, "Captcha not found")
			ctx.Abort()
			return
		}

		if ctx.Request.FormValue("reload") != "" {
			captcha.Reload(id)
		}

		lang := strings.ToLower(ctx.Request.FormValue("lang"))
		var content bytes.Buffer
		switch ext {
		case ".png":
			ctx.Header("Content-Type", "image/png")

			// 宽高
			width := captcha.StdWidth
			height := captcha.StdHeight

			if r.Configuration.CaptchaImageWidth > 0 || r.Configuration.CaptchaImageHeight > 0 {
				width = r.Configuration.CaptchaImageWidth
				height = r.Configuration.CaptchaImageHeight
			}

			captcha.WriteImage(&content, id, width, height)

		case ".wav":
			ctx.Header("Content-Type", "audio/x-wav")
			captcha.WriteAudio(&content, id, lang)

		default:
			ctx.String(http.StatusNotFound, "Captcha not found")
			ctx.Abort()
			return
		}
		http.ServeContent(ctx.Writer, ctx.Request, id+ext, time.Time{}, bytes.NewReader(content.Bytes()))
	})
}
