package msadmin

import (
	"os"
	"math/rand"
	"time"
	"net/http"
	"context"

	"github.com/doublemo/msadmin/config"
	"github.com/doublemo/msadmin/service"
	"github.com/doublemo/msadmin/core/utils"
	"github.com/go-kit/kit/log"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/sd/etcdv3"
)

// Admin 主入口
type Admin struct{
	// registry 全局数据
	r *config.Registry

	// srv http server
	srv *http.Server
}

func (admin *Admin) Serve() error {
	var (
		r      = admin.r
		c      = r.Configuration
		logger = r.Logger
	)

	rand.Seed(time.Now().UnixNano())

	// init etcd
	utils.Assert(admin.etcdv3Client())

	// 注册所有服务
	registrar := etcdv3.NewRegistrar(r.Etcdv3Client, etcdv3.Service{
		Key:   service.MakeKey(c.ETCDFrefix, c.LocalIP+c.GRPCListen),
		Value: c.LocalIP + c.GRPCListen,
	}, logger)

	registrar.Register()
	defer registrar.Deregister()

	// set mode
	gin.SetMode(gin.ReleaseMode)
	if c.Runmode == "debug" {
		gin.SetMode(gin.DebugMode)
	}

	return admin.serveOfGin()
}

func (admin *Admin) serveOfGin() (err error) {
	var (
		r      = admin.r
		c      = r.Configuration
	)

	// middleware
	utils.Assert(admin.middlewares())

	// router
	utils.Assert(admin.routes())

	admin.srv = &http.Server{
		Addr:           c.HTTPListen,
		Handler:        r.Gin,
		ReadTimeout:    c.HTTPReadTimeout,
		WriteTimeout:   c.HTTPWriteTimeout,
		MaxHeaderBytes: c.HTTPMaxHeaderBytes,
	}

	if c.SSL {
		err = admin.srv.ListenAndServeTLS(c.SSLCert, c.SSLKey)
	} else {
		err = admin.srv.ListenAndServe()
	}

	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func (admin *Admin) middlewares() error {
	var (
		r = admin.r
	)

	r.Gin.Use(gin.Logger())
	r.Gin.Use(gin.Recovery())

	// session
	middlewareSession(r)

	// csrf
	middlewareCSRF(r)
	return nil
}

func (admin *Admin) routes() error {
	var (
		r = admin.r
		c = r.Configuration
	)

	// 设置静态文件目录
	r.Gin.Static("/assets", "./resources/assets")
	r.Gin.Static("/media", c.HTTPMediaDir)

	// 加载模板目录
	r.Gin.LoadHTMLGlob("resources/view/*")

	// 系统路由
	routes(r)
	return loadRoutes(r)
}

func (admin *Admin) etcdv3Client() error {
	var (
		c = admin.r.Configuration
		err error
	)

	ops := etcdv3.ClientOptions{
		CACert:        c.ETCDCACert,
		Cert:          c.ETCDCert,
		Key:           c.ETCDKey,
		Username:      c.ETCDUsername,
		Password:      c.ETCDPassword,
		DialTimeout:   c.ETCDDialTimeout,
		DialKeepAlive: c.ETCDDialKeepAlive,
	}

	admin.r.Etcdv3Client, err = etcdv3.NewClient(context.Background(), c.ETCDAddress, ops)
	return err
}

func (admin *Admin) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := admin.srv.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		utils.Assert(err)
	}
	
	select {
	case <-ctx.Done():
		admin.r.Logger.Log("Admin shutdown", "[TIMEOUT]")
	}
}

func New(c *config.Configuration, logger log.Logger) *Admin {
	if logger == nil {
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.WithPrefix(logger, "o", "MSAdmin")
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// 生成全局变量
	var r *config.Registry
	{
		r = config.NewRegistry(c, logger)
	}

	return &Admin{
		r:        r,
	}
}