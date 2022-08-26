package auth

import (
	"context"
	authpb "github.com/iralance/scancar/server/auth/api/gen/v1"
	"github.com/iralance/scancar/server/auth/dao"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type Service struct {
	Mongo          *dao.Mongo
	OpenIDResolver OpenIDResolver
	Logger         *zap.Logger
	TokenGenerator TokenGenerator //生成器
	TokenExpire    time.Duration  //过期时间
	*authpb.UnimplementedAuthServiceServer
}

type OpenIDResolver interface {
	Resolve(code string) (string, error)
}

type TokenGenerator interface {
	GenerateToken(accountId string, expire time.Duration) (string, error)
}

func (s Service) Login(ctx context.Context, request *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	openID, err := s.OpenIDResolver.Resolve(request.Code)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable,
			"cannot resolve openid: %v", err)
	}
	s.Logger.Info("received code", zap.String("code", request.Code))

	accountID, err := s.Mongo.ResolveAccountID(ctx, openID)
	s.Logger.Info(accountID, zap.Error(err))
	if err != nil {
		s.Logger.Error("cannot resolve account id", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}
	tkn, err := s.TokenGenerator.GenerateToken(accountID, s.TokenExpire)
	if err != nil {
		s.Logger.Error("cannot generate token", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}
	return &authpb.LoginResponse{
		AccessToken: tkn,
		ExpiresIn:   int32(s.TokenExpire.Seconds()),
	}, nil
}
