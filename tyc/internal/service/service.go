package service

import (
	"github.com/go-redis/redis"
	"github.com/hatlonely/go-kit/logger"
	"github.com/jinzhu/gorm"
	"time"
)

type TycService struct {
	mysqlCli    *gorm.DB
	redisCli    *redis.Client
	logger 		*logger.Logger
	tokenExpiration time.Duration
	tokenMaxRequest int64
}

func NewTycService(mysqlCli *gorm.DB, redisCli *redis.Client, logger *logger.Logger, opts ...TycServiceOption) (*TycService, error) {
	options := defaultAccountServiceOptions
	for _, opt := range opts {
		opt(&options)
	}

	return &TycService{
		mysqlCli: 	mysqlCli,
		redisCli: 	redisCli,
		logger: 	logger,
		tokenExpiration: options.TokenExpiration,
		tokenMaxRequest: options.TokenMaxRequest,
	}, nil
}

type TycServiceOptions struct {
	TokenExpiration time.Duration
	TokenMaxRequest int64
}

var defaultAccountServiceOptions = TycServiceOptions{
	TokenExpiration: 24 * time.Hour,
	TokenMaxRequest: 300,
}

type TycServiceOption func(options *TycServiceOptions)

func WithTokenExpiration(tokenExpiration time.Duration) TycServiceOption {
	return func(options *TycServiceOptions) {
		options.TokenExpiration = tokenExpiration
	}
}

func WithTokenMaxRequest(tokenMaxRequest int64) TycServiceOption {
	return func(options *TycServiceOptions) {
		options.TokenMaxRequest = tokenMaxRequest
	}
}
