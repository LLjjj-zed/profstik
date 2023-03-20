package webHandler

import (
	"context"
	"github.com/132982317/profstik/cmd/api/rpc"
	kitex "github.com/132982317/profstik/kitex_gen/favorite"
	"github.com/132982317/profstik/pkg/errno"
	"github.com/132982317/profstik/pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
	"strconv"
)

func FavoriteAction(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token")
	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if err != nil || (actionType != 1 && actionType != 2) {
		response.FavoriteListErr(c, errno.ParamErr)
		return
	}
	vid, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		response.FavoriteListErr(c, errno.ParamErr)
		return
	}

	_, err = rpc.FavoriteAction(ctx, &kitex.FavoriteActionRequest{
		Token:      token,
		VideoId:    vid,
		ActionType: int32(actionType),
	})
	if err != nil {
		response.FavoriteListErr(c, errno.RpcConnectErr)
		return
	}
	response.FavoriteOk(c, errno.Success)
}

func FavoriteList(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token")
	uid, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		response.FavoriteListErr(c, errno.ParamErr)
		return
	}

	favoriteList, err := rpc.FavoriteList(ctx, &kitex.FavoriteListRequest{
		UserId: uid,
		Token:  token,
	})
	if err != nil {
		response.FavoriteListErr(c, errno.RpcConnectErr)
		return
	}
	response.FavoriteListOk(c, errno.Success, favoriteList.VideoList)
}
