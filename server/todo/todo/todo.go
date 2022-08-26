package todo

import (
	"context"
	"github.com/iralance/scancar/server/shared/auth"
	todopb "github.com/iralance/scancar/server/todo/api/gen/v1"
	"go.uber.org/zap"
)

type Service struct {
	Logger *zap.Logger
	*todopb.UnimplementedTodoServiceServer
}

func (s Service) CreateTodo(ctx context.Context, request *todopb.CreateTodoRequest) (*todopb.CreateTodoResponse, error) {
	//get accountID from context
	accountId, err := auth.AcountIDFromContext(ctx)
	if err != nil {
		s.Logger.Info("create todo", zap.Error(err))
	}
	s.Logger.Info("create trip", zap.String("title", request.Title), zap.String("account_id", accountId.String()))
	return &todopb.CreateTodoResponse{Word: "title accountID" + accountId.String()}, nil
}
