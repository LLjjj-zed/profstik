package main

import (
	"context"
	"github.com/132982317/profstik/kitex_gen/favorite"
)

type FavoriteServiceImpl struct{}

func (f *FavoriteServiceImpl) FavoriteAction(ctx context.Context, req *favorite.FavoriteActionRequest) (*favorite.FavoriteActionResponse, error) {
	return nil, nil
}

func (f *FavoriteServiceImpl) FavoriteList(ctx context.Context, req *favorite.FavoriteListRequest) (*favorite.FavoriteListResponse, error) {
	return nil, nil
}
