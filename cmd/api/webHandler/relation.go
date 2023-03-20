package webHandler

import (
	"context"
	"github.com/132982317/profstik/cmd/api/rpc"
	kitex "github.com/132982317/profstik/kitex_gen/relation"
	"github.com/132982317/profstik/pkg/errno"
	"github.com/132982317/profstik/pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
	"strconv"
)

func FriendList(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token")
	uid, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		response.FriendListErr(c, errno.ParamErr)
		return
	}
	// 调用rpc
	relationFriendList, err := rpc.RelationFriendList(ctx, &kitex.RelationFriendListRequest{
		UserId: uid,
		Token:  token,
	})
	if err != nil {
		response.FriendListErr(c, errno.RpcConnectErr)
		return
	}
	response.FriendListOk(c, errno.Success, relationFriendList.UserList)
}

func FollowerList(ctx context.Context, c *app.RequestContext) {
	uid, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		response.FollowerListErr(c, errno.ParamErr)
		return
	}
	token := c.Query("token")

	relationFollowerList, err := rpc.RelationFollowerList(ctx, &kitex.RelationFollowerListRequest{
		UserId: uid,
		Token:  token,
	})
	if err != nil {
		response.FollowerListErr(c, errno.RpcConnectErr)
		return
	}
	response.FollowerListOk(c, errno.Success, relationFollowerList.UserList)
}

func FollowList(ctx context.Context, c *app.RequestContext) {
	uid, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		response.FollowListErr(c, errno.ParamErr)
		return
	}
	token := c.Query("token")
	relationFollowList, err := rpc.RelationFollowList(ctx, &kitex.RelationFollowListRequest{
		UserId: uid,
		Token:  token,
	})
	if err != nil {
		response.FollowListErr(c, errno.RpcConnectErr)
		return
	}
	response.FollowListOk(c, errno.Success, relationFollowList.UserList)
}

func RelationAction(ctx context.Context, c *app.RequestContext) {
	tid, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		response.RelationErr(c, errno.ParamErr)
		return
	}
	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if err != nil || (actionType != 1 && actionType != 2) {
		response.RelationErr(c, errno.ParamErr)
		return
	}
	token := c.Query("token")
	if token == "" {
		response.RelationErr(c, errno.ParamErr)
		return
	}

	_, err = rpc.RelationAction(ctx, &kitex.RelationActionRequest{
		Token:      token,
		ToUserId:   tid,
		ActionType: int32(actionType),
	})
	if err != nil {
		response.RelationErr(c, errno.RpcConnectErr)
		return
	}
	response.RelationOk(c, errno.Success)
}
