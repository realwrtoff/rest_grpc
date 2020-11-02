package service

import (
	"context"
	"github.com/hatlonely/go-kit/rpcx"
	sycmApi "github.com/realwrtoff/rest_grpc/sycm/api/gen/go/api"
)

func (s *SycmService) SetCookie(ctx context.Context, req *sycmApi.CookieReq) (*sycmApi.CookieRes, error){
	requestID := rpcx.MetaDataGetRequestID(ctx)

	searchRes := &sycmApi.CookieRes{
		Status: 200,
		Message: "success",
		Data: nil,
	}

	searchRes.Data = requestID
	return searchRes, nil
}