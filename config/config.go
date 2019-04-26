package config

import (
	"time"
)

type Configuration struct {
	// ID 当前服务的唯一标识
	ID string

	// Runmode 支持模式
	Runmode string

	// IDAddress 当前服务器IP地址
	LocalIP string

	// GRPCListen GRPC服务 监听端口
	GRPCListen string

	// HTTPListen HTTP服务 监听端口
	HTTPListen string

	// HTTPReadTimeout http服务读取超时
	HTTPReadTimeout time.Duration

	// HTTPWriteTimeout http服务写入超时
	HTTPWriteTimeout time.Duration

	// HTTPMaxHeaderBytes  http内容大小限制
	HTTPMaxHeaderBytes int

	// HTTPMediaDir 媒体文件目录
	HTTPMediaDir string

	// SSL ssl 支持
	SSL bool

	// SSLKey 证书key
	SSLKey string

	// SSLCert 证书
	SSLCert string

	// SessionName session名称
	SessionName string

	// SessionSecret session验证key
	SessionSecret string

	// CSRFSecret csrf
	CSRFSecret string

	// ETCDAddress etcd 服务器地址
	ETCDAddress []string

	// ETCDFrefix etcd 存储值前缀
	ETCDFrefix string

	// ETCDCACert etcd CA证书地址
	ETCDCACert string

	// ETCDCert etcd 证书地址
	ETCDCert string

	// ETCDCert  etcd 证书key地址
	ETCDKey string

	// ETCDUsername etcd 验证用户名
	ETCDUsername string

	// ETCDPassword etcd 验证密码
	ETCDPassword string

	// ETCDDialTimeout if DialTimeout is 0, it defaults to 3s
	ETCDDialTimeout time.Duration

	//  ETCDDialKeepAlive If DialKeepAlive is 0, it defaults to 3s
	ETCDDialKeepAlive time.Duration

	// MongodbAddress mongodb连接信息
	MongodbAddress string

	// MongodbName 数据库名称
	MongodbName string

	// ServiceSecurityKey JWT 服务之通信认证
	ServiceSecurityKey []byte
}

// NewConfig 创建配置文件
// 默认填充相关数据,在不修改的情况将使用默认数据来运行网关
func New() *Configuration {
	return &Configuration{
		ID:                 "micro-admin",
		Runmode:            "debug",
		LocalIP:            "127.0.0.1",
		GRPCListen:         ":59000",
		HTTPListen:         ":8088",
		ETCDAddress:        []string{"127.0.0.1:2379"},
		ETCDFrefix:         "/services/foxchat",
		ETCDDialTimeout:    3 * time.Second,
		ETCDDialKeepAlive:  3 * time.Second,
		MongodbAddress:     "mongodb://192.168.2.19:27017",
		MongodbName:        "foxchat",
		ServiceSecurityKey: []byte("foxchat"),
		HTTPReadTimeout:    10 * time.Second,
		HTTPWriteTimeout:   10 * time.Second,
		HTTPMaxHeaderBytes: 1 << 20,
		HTTPMediaDir:       "media",
		SessionSecret:      "micro-admin-session-secret",
		SessionName:        "microadmin",
		CSRFSecret:         "micro-admin-csrf-secret",
	}
}