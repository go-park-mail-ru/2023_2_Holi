package grpc

import (
	"2023_2_Holi/domain"
	"2023_2_Holi/domain/grpc/session"
	"context"
	"google.golang.org/grpc"
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
	userID, err := h.AuthUsecase.IsAuth(st.Token)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid session token")
	}
	
	return nil, grpc.Errorf(codes.NotFound, "session not found")

	return nil, status.Errorf(codes.Unimplemented, "method IsAuth not implemented")
}
