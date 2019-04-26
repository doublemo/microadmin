package admin

import (
	"time"
	"net/http"

	"github.com/doublemo/msadmin/config"
	"github.com/doublemo/msadmin/flash"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	csrf "github.com/utrack/gin-csrf"
)

func login(ctx *gin.Context, r *config.Registry) {
	flashes := flash.NewFast()
	session := sessions.Default(ctx)

	var paswordAttemptCounter int
	counter, _ := session.Get("paswordAttemptCounter").(int)
	if counter >= 5 || counter < 1 {
		paswordAttemptCounter = 0
	} else {
		paswordAttemptCounter = 5 - counter
	}

	ctx.HTML(http.StatusOK, "login.tmpl", gin.H{
		"title":                 "Login",
		"csrf":                  csrf.GetToken(ctx),
		"flashes":               flashes,
		"username":              "",
		"paswordAttemptCounter": paswordAttemptCounter,
	})
}



func logout(ctx *gin.Context, r *config.Registry) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
	ctx.Redirect(http.StatusFound, "/login")
}

func postLogin(ctx *gin.Context, r *config.Registry) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	flashes := flash.NewFast()
	h := gin.H{
		"title":                 "Login",
		"csrf":                  csrf.GetToken(ctx),
		"username":              username,
		"flashes":               flashes,
		"paswordAttemptCounter": 0,
	}

	session := sessions.Default(ctx)
	incorrectTime, ok := session.Get("incorrecttime").(int64)
	if ok && incorrectTime > 0 && time.Now().Before(time.Unix(incorrectTime, 0)) {
		flashes.Error("too many incorrect password attempts, please wait 20 minutes and attempt to sign in again.")
		ctx.HTML(http.StatusOK, "login.tmpl", h)
		return
	}

	if username == "" || password == "" {
		flashes.Error("username or password is incorrect. please try again")
	}

	if username != "admin" && password != "123456" {
		var paswordAttemptCounter int
		paswordAttemptCounter, _ = session.Get("paswordAttemptCounter").(int)
		if paswordAttemptCounter+1 >= 5 {
			session.Set("incorrecttime", time.Now().Add(20*time.Minute).Unix())
			h["paswordAttemptCounter"] = 0
		} else {
			h["paswordAttemptCounter"] = 5 - (paswordAttemptCounter + 1)
		}

		session.Set("paswordAttemptCounter", paswordAttemptCounter+1)
		session.Save()
		flashes.Error("username or password is incorrect. please try again")
	}

	if flashes.IsError() {
		ctx.HTML(http.StatusOK, "login.tmpl", h)
		return
	}

	session.Set("UID", uint64(1))
	session.Delete("incorrecttime")
	session.Delete("paswordAttemptCounter")
	session.Save()
	ctx.Redirect(http.StatusMovedPermanently, "/admin")
}