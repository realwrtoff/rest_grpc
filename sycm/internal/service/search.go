package service

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/hatlonely/go-kit/rpcx"
	"github.com/pkg/errors"
	sycm "github.com/realwrtoff/rest_grpc/sycm/api/gen/go/api"
	"github.com/realwrtoff/rest_grpc/sycm/internal/public"
	"google.golang.org/grpc/codes"
	"net/url"
)

func (s *SycmService) SetCookie(ctx context.Context, req *sycm.CookieReq) (*sycm.CookieRes, error) {
	requestID := rpcx.MetaDataGetRequestID(ctx)
	searchRes := &sycm.CookieRes{
		Code:    200,
		Message: "success",
	}
	// token 校验
	cookieMap := public.CookieToMap(req.Cookie)
	s.logger.Info(cookieMap)
	if _, ok := cookieMap["sn"]; !ok {
		searchRes.Code = 404
		searchRes.Message = "account not found in cookie"
		return searchRes, rpcx.NewErrorWithoutReferf(codes.NotFound, requestID, "NotFound", "account not found in cookie")
	}
	var account string
	if public.IsHan(cookieMap["sn"]) {
		account = cookieMap["sn"]
	} else {
		account, _ = url.QueryUnescape(cookieMap["sn"])
	}
	if account != s.options.CookieHashField {
		s.redisCli.HSet(s.options.CookieHashKey, account, req.Cookie)
	}
	if _, err := s.redisCli.HSet(s.options.CookieHashKey, s.options.CookieHashField, req.Cookie).Result(); err != nil {
		searchRes.Code = 500
		searchRes.Message = "set account cookie failed"
		searchRes.Data = err.Error()
		return searchRes, rpcx.NewErrorWithoutReferf(codes.Internal, requestID, "Internal Error", "account[%v] set cookie failed", account)
	}
	// 写入队列，因为发现cookie要过一段时间才能生效
	length, err := s.redisCli.RPush(s.options.CookieQueue, req.Cookie).Result()
	if err != nil {
		searchRes.Code = 500
		searchRes.Message = "rpush account cookie failed"
		searchRes.Data = err.Error()
		return searchRes, errors.Wrapf(err, "rpush cookie %s failed", req.Cookie)
	}
	searchRes.Data = fmt.Sprintf("set %v cookie ok, total %v cookie in queue", account, length)
	return searchRes, nil
}

func (s *SycmService) GetCookie(ctx context.Context, req *sycm.CookieReq) (*sycm.CookieRes, error) {
	requestID := rpcx.MetaDataGetRequestID(ctx)
	searchRes := &sycm.CookieRes{
		Code:    200,
		Message: "success",
	}
	// token 校验
	cookie, err := s.redisCli.HGet(s.options.CookieHashKey, req.Account).Result()
	if err != nil && err != redis.Nil {
		searchRes.Code = 500
		searchRes.Message = "get account cookie failed"
		searchRes.Data = err.Error()
		return searchRes, rpcx.NewErrorWithoutReferf(codes.Internal, requestID, "Internal Error", "account[%v] get cookie failed", req.Account)
	} else if err == redis.Nil && req.Account != s.options.CookieHashField {
		cookie, _ = s.redisCli.HGet(s.options.CookieHashKey, s.options.CookieHashField).Result()
	}
	searchRes.Data = cookie
	return searchRes, nil
}
