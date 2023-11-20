package grpc

import (
	"2023_2_Holi/domain"
	"2023_2_Holi/domain/grpc/session"
	logs "2023_2_Holi/logger"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	session.UnimplementedAuthCheckerServer
	AuthUsecase domain.AuthUsecase
}

func NewAuthAuthHandler(u domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		AuthUsecase: u,
	}
}

func (h *AuthHandler) IsAuth(ctx context.Context, st *session.Token) (*session.UserID, error) {
	if st.Token == "" {
		return nil, status.Errorf(codes.InvalidArgument, "invalid session token")
	}
	logs.Logger.Debug("IsAuth token:", st.Token)

	userID, err := h.AuthUsecase.IsAuth(st.Token)
	if err != nil {
		return nil, status.Errorf(domain.GetGrpcStatusCode(err), err.Error())
	}
	logs.Logger.Debug("IsAuth userID:", userID)

	return &session.UserID{
		ID: userID,
	}, nil
}
