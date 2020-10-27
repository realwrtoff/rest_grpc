package service

import (
	"context"
	"encoding/hex"
	"regexp"

	"github.com/hatlonely/go-kit/rpcx"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"

	tyc "github.com/realwrtoff/rest_grpc/tyc/api/gen/go/api"
	"github.com/realwrtoff/rest_grpc/tyc/internal/model"
)

func GenerateToken() string {
	return hex.EncodeToString(uuid.NewV4().Bytes())
}

var ReUserName = regexp.MustCompile(`^[a-zA-Z0-9_]\w+$`)

func (s *TycService) SignIn(ctx context.Context, req *tyc.SignInReq) (*tyc.SignInRes, error) {
	requestID := rpcx.MetaDataGetRequestID(ctx)

	signInRes := &tyc.SignInRes{
		Status: 200,
		Message: "success",
		Token: "",
	}

	a := &model.Account{}

	// 如何防止sql注入
	if ! ReUserName.MatchString(req.Username) {
		return signInRes, rpcx.NewErrorWithoutReferf(codes.InvalidArgument, requestID, "InvalidArgument", "user [%v] is invalid", req.Username)
	}

	// 查找用户
	if err := s.mysqlCli.Where("username=?", req.Username).First(a).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return signInRes, rpcx.NewErrorWithoutReferf(codes.NotFound, requestID, "Forbidden", "user [%v] not exist", req.Username)
		}
		return signInRes, errors.Wrapf(err, "mysql select user [%v] failed", req.Username)
	}

	if a.Password != req.Password {
		return nil, rpcx.NewErrorWithoutReferf(codes.PermissionDenied, requestID, "Forbidden", "password is incorrect")
	}

	// 生成token
	token := GenerateToken()
	userToken := &model.Token{
		Username: req.Username,
		Token: token,
	}
	// 记录token
	if err := s.mysqlCli.Create(userToken).Error; err != nil {
		return signInRes, rpcx.NewErrorWithoutReferf(codes.Internal, requestID, "InternalServerError", "insert into mysql failed [%v]", err.Error())
	}

	// 写入redis, 默认有效期1天
	if res, err := s.redisCli.Set(token, 0, s.tokenExpiration).Result(); err != nil {
		return signInRes, errors.Wrapf(err, "redis set token [%s] failed, res[%s]", token, res)
	}
	signInRes.Token = token
	return signInRes, nil
}
