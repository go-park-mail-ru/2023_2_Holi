package grpc

import (
	"2023_2_Holi/domain"
	"2023_2_Holi/domain/grpc/subscription"
	logs "2023_2_Holi/logger"
	"context"
	"google.golang.org/grpc/status"
	"strconv"
)

type SubHandler struct {
	subscription.UnimplementedSubCheckerServer
	SubsUsecase domain.SubsUsecase
}

func NewSubHandler(u domain.SubsUsecase) *SubHandler {
	return &SubHandler{
		SubsUsecase: u,
	}
}

func (h *SubHandler) CheckSub(ctx context.Context, userID *subscription.UserID) (*subscription.Result, error) {
	uid, err := strconv.Atoi(userID.ID)
	subUpTo, checkStatus, err := h.SubsUsecase.CheckSub(uid)
	if err != nil {
		return nil, status.Errorf(domain.GetGrpcStatusCode(err), err.Error())
	}
	logs.Logger.Info("CheckSub subupTo:", subUpTo)

	return &subscription.Result{
		Status: checkStatus,
	}, nil
}
