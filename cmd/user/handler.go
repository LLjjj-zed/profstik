package main

import (
	"context"
	"github.com/132982317/profstik/kitex_gen/user"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

func (u *UserServiceImpl) Register(ctx context.Context, req *user.UserRegisterRequest) (*user.UserRegisterResponse, error) {
	return nil, nil
}

func (u *UserServiceImpl) Login(ctx context.Context, req *user.UserLoginRequest) (*user.UserLoginResponse, error) {
	return nil, nil
}

func (u *UserServiceImpl) UserInfo(ctx context.Context, req *user.UserInfoRequest) (*user.UserInfoResponse, error) {
	return nil, nil
}
