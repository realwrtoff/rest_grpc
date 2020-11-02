package service

import (
	"context"
	"github.com/hatlonely/go-kit/rpcx"
	sycm "github.com/realwrtoff/rest_grpc/sycm/api/gen/go/api"
	"google.golang.org/grpc/codes"
	"strings"
)

func CookieToMap(cookieStr string) map[string]string  {
	cookies := strings.Split(cookieStr, "; ")
	mp := make(map[string]string, len(cookies))
	for _, cookie := range cookies {
		kvs := strings.Split(cookie, "=")
		mp[kvs[0]] = kvs[1]
	}
	return mp
}

func (s *SycmService) SetCookie(ctx context.Context, req *sycm.CookieReq) (*sycm.CookieRes, error){
	requestID := rpcx.MetaDataGetRequestID(ctx)
	searchRes := &sycm.CookieRes{
		Code: 200,
		Message: "success",
	}
	// token 校验
	cookieMap := CookieToMap(req.Cookie)
	if _, ok := cookieMap["sn"]; !ok {
		searchRes.Code = 404
		searchRes.Message = "account not found in cookie"
		return searchRes, rpcx.NewErrorWithoutReferf(codes.NotFound, requestID, "NotFound", "account not found in cookie")
	}
	if _, err := s.redisCli.HSet(s.cookieHashKey, cookieMap["sn"], req.Cookie).Result(); err != nil {
		searchRes.Code = 500
		searchRes.Message = "set account cookie failed"
		searchRes.Data = err.Error()
		return searchRes, rpcx.NewErrorWithoutReferf(codes.Internal, requestID, "Internal Error", "account[%v] set cookie failed", cookieMap["sn"])
	}
	searchRes.Data = req.Cookie
	return searchRes, nil
}

func (s *SycmService) GetCookie(ctx context.Context, req *sycm.CookieReq) (*sycm.CookieRes, error){
	requestID := rpcx.MetaDataGetRequestID(ctx)
	searchRes := &sycm.CookieRes{
		Code: 200,
		Message: "success",
	}
	// token 校验
	cookie, err := s.redisCli.HGet(s.cookieHashKey, req.Account).Result()
	if err != nil {
		searchRes.Code = 500
		searchRes.Message = "get account cookie failed"
		searchRes.Data = err.Error()
		return searchRes, rpcx.NewErrorWithoutReferf(codes.Internal, requestID, "Internal Error", "account[%v] get cookie failed", req.Account)
	}
	searchRes.Data = cookie
	return searchRes, nil
}