package response

import (
	"github.com/132982317/profstik/kitex_gen/message"
	"github.com/132982317/profstik/pkg/errno"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

type MessageChatResponse struct {
	common
	MessageList []*message.Message `json:"message_list"`
}

// MessageChatOk 返回正确信息
func MessageChatOk(c *app.RequestContext, err error, MessageList []*message.Message) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, &MessageChatResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
		MessageList: MessageList,
	})
}

// MessageChatErr  返回错误信息
func MessageChatErr(c *app.RequestContext, err error) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, &MessageChatResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
	})
}

type MessageResponse struct {
	common
}

// MessageOk 返回正确信息
func MessageOk(c *app.RequestContext, err error) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, &MessageResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
	})
}

// MessageErr  返回错误信息
func MessageErr(c *app.RequestContext, err error) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, &MessageResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
	})
}
