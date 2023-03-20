package webHandler

import (
	"context"
	"github.com/132982317/profstik/cmd/api/rpc"
	kitex "github.com/132982317/profstik/kitex_gen/message"
	"github.com/132982317/profstik/pkg/errno"
	"github.com/132982317/profstik/pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
	"strconv"
)

func MessageChat(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token")
	toUserID, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		response.MessageChatErr(c, errno.ParamErr)
		return
	}

	// 调用rpc
	messageChat, err := rpc.MessageChat(ctx, &kitex.MessageChatRequest{
		Token:    token,
		ToUserId: toUserID,
	})
	if err != nil {
		response.MessageChatErr(c, errno.RpcConnectErr)
		return
	}
	response.MessageChatOk(c, errno.Success, messageChat.MessageList)
}

func MessageAction(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token")
	if token == "" {
		response.MessageErr(c, errno.ParamErr)
		return
	}

	toUserID, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		response.MessageErr(c, errno.ParamErr)
		return
	}
	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if err != nil || actionType != 1 {
		response.MessageErr(c, errno.ParamErr)
		return
	}

	if len(c.Query("content")) == 0 {
		response.MessageErr(c, errno.ParamErr)
		return
	}

	// 调用rpc
	_, err = rpc.MessageAction(ctx, &kitex.MessageActionRequest{
		Token:      token,
		ToUserId:   toUserID,
		ActionType: int32(actionType),
		Content:    c.Query("content"),
	})
	if err != nil {
		response.MessageErr(c, errno.RpcConnectErr)
		return
	}
	response.MessageOk(c, errno.Success)
}
