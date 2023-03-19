package response

import (
	"github.com/132982317/profstik/kitex_gen/video"
	"github.com/132982317/profstik/pkg/errno"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

type PublishResponse struct {
	common
}

// PublishOk 返回正确信息
func PublishOk(c *app.RequestContext, err error) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, &PublishResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
	})
}

// PublishErr  返回错误信息
func PublishErr(c *app.RequestContext, err error) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, &PublishResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
	})
}

type PublishListResponse struct {
	common
	VideoList []*video.Video `json:"video_list"`
}

// PublishListOK 返回正确信息
func PublishListOK(c *app.RequestContext, err error, list []*video.Video) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, &PublishListResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
		VideoList: list,
	})
}

// PublishListErr 返回错误信息
func PublishListErr(c *app.RequestContext, err error) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, &PublishListResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
	})
}

type FeedResponse struct {
	common
	NextTime  int64          `json:"next_time"`
	VideoList []*video.Video `json:"video_list"`
}

// FeedOK 返回正确信息
func FeedOK(c *app.RequestContext, err error, videos []*video.Video, nextTime int64) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, FeedResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
		NextTime:  nextTime,
		VideoList: videos,
	})
}

// FeedErr 返回错误信息
func FeedErr(c *app.RequestContext, err error) {
	Err := errno.ConvertErr(err)
	c.JSON(201, FeedResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
	})
}
