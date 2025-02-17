package service

import (
	"context"

	"github.com/mstgnz/starter-kit/api/model"
)

type userService struct {
}

func NewUserService() *userService {
	return &userService{}
}

func (s *userService) Login(ctx context.Context, login *model.Login) (*model.User, error) {
	return nil, nil
}
