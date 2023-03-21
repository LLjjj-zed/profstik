package response

import (
	"context"
	"github.com/132982317/profstik/pkg/errno"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

type common struct {
	StatusCode int64  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

func ResponseWithError(ctx context.Context, c *app.RequestContext, request int, s string) {
	c.JSON(http.StatusOK, &FavoriteResponse{
		common: common{
			StatusCode: errno.ServiceErrCode,
			StatusMsg:  errno.ServiceErr.Error(),
		},
	})
}
