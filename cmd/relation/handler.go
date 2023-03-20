package main

import (
	"context"
	"github.com/132982317/profstik/kitex_gen/relation"
)

type RelationServiceImpl struct{}

func (r *RelationServiceImpl) RelationAction(ctx context.Context, req *relation.RelationActionRequest) (*relation.RelationActionResponse, error) {

}

func (r *RelationServiceImpl) RelationFollowList(ctx context.Context, req *relation.RelationFollowListRequest) (*relation.RelationFollowListResponse, error) {

}

func (r *RelationServiceImpl) RelationFollowerList(ctx context.Context, req *relation.RelationFollowerListRequest) (*relation.RelationFollowerListResponse, error) {

}

func (r *RelationServiceImpl) RelationFriendList(ctx context.Context, req *relation.RelationFriendListRequest) (*relation.RelationFriendListResponse, error) {

}
