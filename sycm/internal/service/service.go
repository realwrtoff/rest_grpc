package service

import (
	"github.com/go-redis/redis"
	"github.com/hatlonely/go-kit/logger"
	"github.com/jinzhu/gorm"
	"time"
)

type SycmService struct {
	mysqlCli    *gorm.DB
	redisCli    *redis.Client
	logger 		*logger.Logger
	tokenExpiration time.Duration
	tokenMaxRequest int64
}

func NewSycmService(mysqlCli *gorm.DB, redisCli *redis.Client, logger *logger.Logger, opts ...SycmServiceOption) (*SycmService, error) {
	options := defaultAccountServiceOptions
	for _, opt := range opts {
		opt(&options)
	}

	return &SycmService{
		mysqlCli: 	mysqlCli,
		redisCli: 	redisCli,
		logger: 	logger,
		tokenExpiration: options.TokenExpiration,
		tokenMaxRequest: options.TokenMaxRequest,
	}, nil
}

type SycmServiceOptions struct {
	TokenExpiration time.Duration
	TokenMaxRequest int64
}

var defaultAccountServiceOptions = SycmServiceOptions{
	TokenExpiration: 24 * time.Hour,
	TokenMaxRequest: 300,
}

type SycmServiceOption func(options *SycmServiceOptions)

func WithTokenExpiration(tokenExpiration time.Duration) SycmServiceOption {
	return func(options *SycmServiceOptions) {
		options.TokenExpiration = tokenExpiration
	}
}

func WithTokenMaxRequest(tokenMaxRequest int64) SycmServiceOption {
	return func(options *SycmServiceOptions) {
		options.TokenMaxRequest = tokenMaxRequest
	}
}