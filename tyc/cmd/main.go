package main

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	tyc "github.com/realwrtoff/rest_grpc/tyc/api/gen/go/api"
)

type MobileService struct{}

func (s *MobileService) Echo(ctx context.Context, req *tyc.MobileReq) (*tyc.MobileRes, error) {
	return &tyc.MobileRes{CompanyId: req.CompanyId, Name: req.Name, Phone:"17744581949"}, nil
}


func main() {
	mux := runtime.NewServeMux()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if err := tyc.RegisterMobileServiceHandlerServer(ctx, mux, &MobileService{}); err != nil {
		panic(err)
	}
	if err := http.ListenAndServe(":80", mux); err != nil {
		panic(err)
	}
}
