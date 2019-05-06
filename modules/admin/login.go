package admin

import (
	"net/http"
	"time"

	"github.com/doublemo/msadmin/config"
	"github.com/doublemo/msadmin/flash"
	"github.com/doublemo/msadmin/modules/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func login(ctx *gin.Context, r *config.Registry) {
	ctx.HTML(http.StatusOK, "login.tmpl", makeLoginPageData(ctx))
}

func logout(ctx *gin.Context, r *config.Registry) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
	ctx.Redirect(http.StatusFound, "/login")
}

func postLogin(ctx *gin.Context, r *config.Registry) {
	session := sessions.Default(ctx)
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	h := makeLoginPageData(ctx)

	incorrectTime, ok := session.Get("incorrecttime").(int64)
	if ok && chkIncorrectCounter(incorrectTime) {
		h["flashes"].(flash.Fast).Error("too many incorrect password attempts, please wait 20 minutes and attempt to sign in again.")
		ctx.HTML(http.StatusOK, "login.tmpl", h)
		return
	}

	captchaID, ok := session.Get("captchaID").(string)
	if ok && chkCaptcha(h["paswordAttemptCounter"].(int), captchaID, ctx.PostForm("captcha")) {
		h["flashes"].(flash.Fast).Error("Wrong captcha solution!")
		ctx.HTML(http.StatusOK, "login.tmpl", h)
		return
	}

	if username != "admin" && password != "123456" {
		paswordAttemptCounter := h["paswordAttemptCounter"].(int)
		paswordAttemptCounter++
		if paswordAttemptCounter >= 5 {
			session.Set("incorrecttime", time.Now().Add(20*time.Minute).Unix())
		}

		if paswordAttemptCounter > 0 && paswordAttemptCounter%2 == 0 {
			captchaID := captcha.New()
			h["captchaID"] = captchaID
			session.Set("captchaID", captchaID)
		} else if paswordAttemptCounter > 2 {
			h["captchaID"] = ""
			session.Delete("captchaID")
		}

		h["paswordAttemptCounter"] = paswordAttemptCounter
		session.Set("paswordAttemptCounter", paswordAttemptCounter)
		session.Save()
		h["flashes"].(flash.Fast).Error("username or password is incorrect. please try again")
	}

	if h["flashes"].(flash.Fast).IsError() {
		ctx.HTML(http.StatusOK, "login.tmpl", h)
		return
	}

	session.Set("UID", uint64(1))
	session.Delete("captchaID")
	session.Delete("paswordAttemptCounter")
	session.Delete("incorrecttime")
	session.Save()
	ctx.Redirect(http.StatusMovedPermanently, "/admin")
}

func chkIncorrectCounter(incorrectTime int64) bool {
	if incorrectTime > 0 && time.Now().Before(time.Unix(incorrectTime, 0)) {
		return true
	}
	return false
}

func chkCaptcha(incorrectCounter int, captchaID, captchaVal string) bool {
	if incorrectCounter > 1 && !captcha.VerifyString(captchaID, captchaVal) {
		captcha.Reload(captchaID)
		return true
	}

	return false
}

func makeLoginPageData(ctx *gin.Context) gin.H {
	session := sessions.Default(ctx)
	username := ctx.PostForm("username")
	flashes := flash.NewFast()

	var paswordAttemptCounter int
	{
		paswordAttemptCounter, _ = session.Get("paswordAttemptCounter").(int)
	}

	var captchaID string
	{
		captchaID, _ = session.Get("captchaID").(string)
	}

	return gin.H{
		"title":                 "Login",
		"csrf":                  csrf.GetToken(ctx),
		"username":              username,
		"flashes":               flashes,
		"paswordAttemptCounter": paswordAttemptCounter,
		"captchaID":             captchaID,
	}
}
