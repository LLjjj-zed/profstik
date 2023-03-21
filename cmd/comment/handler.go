package main

import (
	"context"
	"github.com/132982317/profstik/kitex_gen/comment"
)

type CommentServiceImpl struct{}

func (c *CommentServiceImpl) CommentAction(ctx context.Context, req *comment.CommentActionRequest) (*comment.CommentActionResponse, error) {
	return nil, nil
}

func (c *CommentServiceImpl) CommentList(ctx context.Context, req *comment.CommentListRequest) (*comment.CommentListResponse, error) {
	return nil, nil
}
