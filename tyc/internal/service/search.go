package service

import (
	"context"
	"encoding/json"
	"github.com/hatlonely/go-kit/rpcx"
	"github.com/pkg/errors"
	tyc "github.com/realwrtoff/rest_grpc/tyc/api/gen/go/api"
	"strconv"
)

func (s *TycService) Search(ctx context.Context, req *tyc.SearchReq) (*tyc.SearchRes, error){
	requestID := rpcx.MetaDataGetRequestID(ctx)

	searchRes := &tyc.SearchRes{
		Status: 0,
		Message: "success",
		Data: nil,
	}

	// 校验token是否存在
	exists, err := s.redisCli.Exists(req.Token).Result()
	if err != nil {
		searchRes.Message = err.Error()
		return searchRes, err
	}
	if exists < 1 {
		searchRes.Message = "token not found"
		return searchRes, errors.Wrapf(err, "token %v not found", req.Token)
	}

	// 校验token请求数量是否达到限制
	cnt, err := s.redisCli.Incr(req.Token).Result()
	if err != nil {
		searchRes.Message = err.Error()
		return searchRes, err
	}
	if cnt > 100 {
		searchRes.Message = "max request limited"
		return searchRes, errors.Wrapf(err, "token %v max request limited", req.Token)
	}

	// 查找company信息
	res, err := s.redisCli.Get(strconv.FormatInt(req.CompanyId,10)).Bytes()
	if err != nil {
		searchRes.Message = err.Error()
		return searchRes, err
	}
	tycMobile := &tyc.TycMobile{}
	if err := json.Unmarshal(res, tycMobile); err != nil {
		searchRes.Message = err.Error()
		return searchRes, err
	}
	// 校验查询结果是否和请求公司名称匹配
	if tycMobile.Name != req.Name {
		searchRes.Message = "request info not matched"
		return searchRes, errors.Wrapf(err, "request[%v] comapnyId[%v]not match with company name[%v]", requestID, req.CompanyId, req.Name)
	}
	searchRes.Data = tycMobile
	return searchRes, nil
}