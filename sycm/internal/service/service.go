package service

import (
	"github.com/go-redis/redis"
	"github.com/hatlonely/go-kit/logger"
	"github.com/jinzhu/gorm"
	"time"
)

type Options struct {
	TokenExpiration time.Duration
	TokenMaxRequest int64
	CookieHashKey   string
	CookieHashField string
	CookieQueue     string
}

type SycmService struct {
	mysqlCli *gorm.DB
	redisCli *redis.Client
	logger   *logger.Logger
	options  *Options
}

func NewSycmService(mysqlCli *gorm.DB, redisCli *redis.Client, logger *logger.Logger, options *Options) (*SycmService, error) {
	return &SycmService{
		mysqlCli: mysqlCli,
		redisCli: redisCli,
		logger:   logger,
		options:  options,
	}, nil
}
