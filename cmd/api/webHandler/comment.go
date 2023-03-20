package webHandler

import (
	"context"
	"github.com/132982317/profstik/cmd/api/rpc"
	kitex "github.com/132982317/profstik/kitex_gen/comment"
	"github.com/132982317/profstik/pkg/errno"
	"github.com/132982317/profstik/pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
	"strconv"
)

func CommentAction(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token")
	vid, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		response.CommentErr(c, errno.ParamErr)
		return
	}
	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if err != nil || (actionType != 1 && actionType != 2) {
		response.CommentErr(c, errno.ParamErr)
		return
	}
	req := new(kitex.CommentActionRequest)
	req.Token = token
	req.VideoId = vid
	req.ActionType = int32(actionType)

	if actionType == 1 {
		commentText := c.Query("comment_text")
		if commentText == "" {
			response.CommentErr(c, errno.ParamErr) //nil
			return
		}
		req.CommentText = commentText
	} else if actionType == 2 {
		commentID, err := strconv.ParseInt(c.Query("comment_id"), 10, 64)
		if err != nil {
			response.CommentErr(c, errno.ParamErr) //unlegal
			return
		}
		req.CommentId = commentID
	}
	comment, err := rpc.CommentAction(ctx, req)
	if err != nil {
		response.CommentErr(c, errno.RpcConnectErr)
		return
	}
	response.CommentOk(c, errno.Success, comment.Comment)
}

func CommentList(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token")
	vid, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		response.CommentListErr(c, errno.ParamErr)
		return
	}
	commentList, err := rpc.CommentList(ctx, &kitex.CommentListRequest{
		Token:   token,
		VideoId: vid,
	})
	if err != nil {
		response.CommentListErr(c, errno.RpcConnectErr)
		return
	}
	response.CommentListOk(c, errno.Success, commentList.CommentList)
}
