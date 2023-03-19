package response

import (
	"github.com/132982317/profstik/kitex_gen/video"
	"github.com/132982317/profstik/pkg/errno"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

type FavoriteResponse struct {
	common
}

// FavoriteOk 返回正确信息
func FavoriteOk(c *app.RequestContext, err error) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, &FavoriteResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
	})
}

// FavoriteErr  返回错误信息
func FavoriteErr(c *app.RequestContext, err error) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, &FavoriteResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
	})
}

type FavoriteListResponse struct {
	common
	VideoList []*video.Video `json:"video_list"`
}

// FavoriteListOk 返回正确信息
func FavoriteListOk(c *app.RequestContext, err error, VideoList []*video.Video) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, &FavoriteListResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
		VideoList: VideoList,
	})
}

// FavoriteListErr  返回错误信息
func FavoriteListErr(c *app.RequestContext, err error) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, &FavoriteListResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
	})
}
