package service

import (
	"context"
	"github.com/hatlonely/go-kit/rpcx"
	sycm "github.com/realwrtoff/rest_grpc/tyc/api/gen/go/api"
)

func (s *SycmService) SetCookie(ctx context.Context, req *sycm.CookieReq) (*sycm.CookieRes, error){
	requestID := rpcx.MetaDataGetRequestID(ctx)

	searchRes := &sycm.CookieRes{
		Status: 200,
		Message: "success",
		Data: nil,
	}

	searchRes.Data = requestID
	return searchRes, nil
}