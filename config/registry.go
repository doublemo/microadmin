package config

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Registry struct {
	// Configuration 相关配置
	Configuration *Configuration

	// Logger 日志处理
	Logger log.Logger

	// Etcdv3Client 连接实例,用于服务发现
	Etcdv3Client etcdv3.Client

	// Gin gin.Engine
	Gin *gin.Engine

	// DB 数据库连接
	DB *gorm.DB
}

func NewRegistry(c *Configuration, logger log.Logger) *Registry {
	return &Registry{
		Configuration: c,
		Logger:        logger,
		Gin:           gin.New(),
	}
}