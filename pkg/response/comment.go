package response

import (
	"github.com/132982317/profstik/kitex_gen/comment"
	"github.com/132982317/profstik/pkg/errno"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

type CommentResponse struct {
	common
	Comment *comment.Comment `json:"comment"`
}

// CommentOk 返回正确信息
func CommentOk(c *app.RequestContext, err error, comment *comment.Comment) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, &CommentResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
		Comment: comment,
	})
}

// CommentErr  返回错误信息
func CommentErr(c *app.RequestContext, err error) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, &CommentResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
	})
}

type CommentListResponse struct {
	common
	CommentList []*comment.Comment `json:"comment_list"`
}

// CommentListOk 返回正确信息
func CommentListOk(c *app.RequestContext, err error, commentList []*comment.Comment) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, &CommentListResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
		CommentList: commentList,
	})
}

// CommentListErr  返回错误信息
func CommentListErr(c *app.RequestContext, err error) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, &CommentListResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
	})
}
