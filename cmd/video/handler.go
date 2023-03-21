package main

import (
	"context"
	"github.com/132982317/profstik/kitex_gen/video"
)

type VideoServiceImpl struct{}

func (v *VideoServiceImpl) Feed(ctx context.Context, req *video.FeedRequest) (*video.FeedResponse, error) {
	return nil, nil
}

func (v *VideoServiceImpl) PublishAction(ctx context.Context, req *video.PublishActionRequest) (*video.PublishActionResponse, error) {
	return nil, nil
}

func (v *VideoServiceImpl) PublishList(ctx context.Context, req *video.PublishListRequest) (*video.PublishListResponse, error) {
	return nil, nil
}
