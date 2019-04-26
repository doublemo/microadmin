package main 

import (
	"os"
	"log"
	"os/signal"
	"syscall"
	
	"github.com/doublemo/msadmin"
	"github.com/doublemo/msadmin/config"
	"github.com/urfave/cli"
)

var (
	// version 编译版本号
	version string

	// commitID  git commit服务器版本号
	commitID string

	// buildAt 编译日期
	buildAt string
)

func main() {
	app := cli.NewApp()
	app.Name = "MicroAdmin"
	app.Usage = "微服务数据管理平台"
	app.Version = version
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "config,c",
			Value:   "microadmin.conf",
			Usage:   "Load configuration from `FILE`",
		},

		&cli.StringFlag{
			Name:    "id",
			Value:   "microadmin",
			Usage:   "ID 当前服务的唯一标识",
		},

		&cli.StringFlag{
			Name:    "mode",
			Value:   "debug",
			Usage:   "支持模式, 支持 debug release",
		},

		&cli.StringFlag{
			Name:    "localip",
			Value:   "127.0.0.1",
			Usage:   "当前服务器IP地址",
		},

		&cli.StringFlag{
			Name:    "grpc.addr",
			Value:   ":8099",
			Usage:   "GRPC服务监听端口",
		},

		&cli.StringFlag{
			Name:    "http.addr",
			Value:   ":8099",
			Usage:   "HTTP服务监听端口",
		},

		&cli.IntFlag{
			Name:    "http.readtimeout",
			Value:   10,
			Usage:   "http服务读取超时(s)",
		},

		&cli.IntFlag{
			Name:    "http.writetimeout",
			Value:   10,
			Usage:   "http服务写入超时(s)",
		},

		&cli.IntFlag{
			Name:    "http.maxheaderbytes ",
			Value:   1 << 20,
			Usage:   "http内容大小限制(B)",
		},

		&cli.StringFlag{
			Name:    "http.assetsdir ",
			Value:   "./assets",
			Usage:   "静态资源目录",
		},

		&cli.StringFlag{
			Name:    "http.templatesdir ",
			Value:   "templates",
			Usage:   "模板文件目录",
		},

		&cli.StringFlag{
			Name:    "http.mediadir ",
			Value:   "./media",
			Usage:   "媒体文件目录",
		},

		&cli.BoolFlag{
			Name:    "http.ssl",
			Usage:   "ssl支持",
		},

		&cli.StringFlag{
			Name:    "http.sslkey",
			Value:   "",
			Usage:   "ssl证书key",
		},

		&cli.StringFlag{
			Name:    "http.sslcert",
			Value:   "",
			Usage:   "ssl证书",
		},

		&cli.StringFlag{
			Name:    "http.sessionsecret",
			Value:   "micro-admin",
			Usage:   "session验证key",
		},

		&cli.StringFlag{
			Name:    "http.csrf",
			Value:   "micro-admin",
			Usage:   "csrf",
		},

		&cli.StringSliceFlag{
			Name:    "etcd.addr",
			Value:   &cli.StringSlice{"127.0.0.1:2379"},
			Usage:   "etcd 服务器地址",
		},

		&cli.StringFlag{
			Name:    "etcd.prefix",
			Value:   "/services/microadmin",
			Usage:   "etcd 存储值前缀",
		},

		&cli.StringFlag{
			Name:    "etcd.cacert",
			Value:   "",
			Usage:   "etcd CA证书地址",
		},

		&cli.StringFlag{
			Name:    "etcd.cert",
			Value:   "",
			Usage:   "etcd 证书地址",
		},

		&cli.StringFlag{
			Name:    "etcd.key",
			Value:   "",
			Usage:   "etcd 证书key地址",
		},

		&cli.StringFlag{
			Name:    "etcd.username",
			Value:   "",
			Usage:   "etcd 验证用户名",
		},

		&cli.StringFlag{
			Name:    "etcd.password",
			Value:   "",
			Usage:   "etcd 验证密码",
		},

		&cli.IntFlag{
			Name:    "etcd.dialTimeout",
			Value:   3,
			Usage:   "etcd 连接超时(s)",
		},

		&cli.IntFlag{
			Name:    "etcd.dialkeepalive",
			Value:   3,
			Usage:   "etcd 超时(s)",
		},

		&cli.StringFlag{
			Name:    "service.securitykey",
			Value:   "micro-admin",
			Usage:   "服务之通信认证",
		},

	}
	app.Action = action
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func action(ctx *cli.Context) error {
	cli.ShowVersion(ctx)
	log.Println("CommitID", commitID)
	log.Println("BuildAt", buildAt)

	c := config.New()
	s := msadmin.New(c, nil)
	go sigHandler(s)
	return s.Serve()
}



func sigHandler(s *msadmin.Admin) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	for {
		msg := <-ch
		switch msg {
		case syscall.SIGINT, syscall.SIGTERM:
			log.Println("Shutdown Server ...")
			s.Shutdown()
			os.Exit(0)
		}
	}
}