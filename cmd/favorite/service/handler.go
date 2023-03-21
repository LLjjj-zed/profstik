package service

import (
	"context"
	"github.com/132982317/profstik/kitex_gen/favorite"
	"github.com/132982317/profstik/pkg/utils/zap"
	Zap "go.uber.org/zap"
)

type FavoriteServiceImpl struct{}

func (f *FavoriteServiceImpl) FavoriteAction(ctx context.Context, req *favorite.FavoriteActionRequest) (resp *favorite.FavoriteActionResponse, err error) {
	logger := zap.LoggerPool.Get().(*Zap.SugaredLogger)
	defer zap.LoggerPool.Put(logger)
	return nil, nil
}

func (f *FavoriteServiceImpl) FavoriteList(ctx context.Context, req *favorite.FavoriteListRequest) (resp *favorite.FavoriteListResponse, err error) {
	logger := zap.LoggerPool.Get().(*Zap.SugaredLogger)
	defer zap.LoggerPool.Put(logger)
	return nil, nil
}
