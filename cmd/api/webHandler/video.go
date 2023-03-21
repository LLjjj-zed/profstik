package webHandler

import (
	"bytes"
	"context"
	"github.com/132982317/profstik/cmd/api/rpc"
	kitex "github.com/132982317/profstik/kitex_gen/video"
	"github.com/132982317/profstik/pkg/errno"
	"github.com/132982317/profstik/pkg/response"
	"github.com/132982317/profstik/pkg/utils/zap"
	"github.com/cloudwego/hertz/pkg/app"
	"io"
	"strconv"
	"time"
)

func Feed(ctx context.Context, c *app.RequestContext) {
	var feedVar FeedParam
	if err := c.Bind(&feedVar); err != nil {
		response.FeedErr(c, errno.ConvertErr(err))
		return
	}
	var timestamp int64 = 0
	if feedVar.LatestTime != "" {
		timestamp, _ = strconv.ParseInt(feedVar.LatestTime, 10, 64)
	} else {
		timestamp = time.Now().UnixMilli()
	}

	feed, err := rpc.Feed(ctx, &kitex.FeedRequest{
		LatestTime: timestamp,
		Token:      feedVar.Token,
	})
	if err != nil {
		response.FeedErr(c, errno.RpcConnectErr)
	}
	if feed.StatusCode == errno.ServiceErrCode {
		response.FeedErr(c, errno.ServiceErr)
		return
	}
	response.FeedOK(c, errno.Success, feed.VideoList, feed.NextTime)
}

func PublishList(ctx context.Context, c *app.RequestContext) {
	token := c.GetString("token")
	uid, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		response.PublishListErr(c, errno.ParamErr)
		return
	}
	publishList, err := rpc.PublishList(ctx, &kitex.PublishListRequest{
		Token:  token,
		UserId: uid,
	})
	if err != nil {
		response.PublishListErr(c, errno.RpcConnectErr)
	}
	if publishList.StatusCode == errno.ServiceErrCode {
		response.PublishListErr(c, errno.ServiceErr)
		return
	}
	response.PublishListOK(c, errno.Success, publishList.VideoList)
}

func PublishAction(ctx context.Context, c *app.RequestContext) {
	logger := zap.InitLogger()
	token := c.PostForm("token")
	if token == "" {
		response.PublishErr(c, errno.AuthorizationFailedErr)
		return
	}
	title := c.PostForm("title")
	if title == "" {
		response.PublishErr(c, errno.ParamErr)
		return
	}
	// 视频数据
	file, err := c.FormFile("data")
	if err != nil {
		logger.Errorln(err.Error())
		response.PublishErr(c, errno.UploadVideoErr)
		return
	}
	src, err := file.Open()
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, src); err != nil {
		logger.Errorln(err.Error())
		response.PublishErr(c, errno.UploadVideoErr)
		return
	}

	pubilsh, err := rpc.PublishAction(ctx, &kitex.PublishActionRequest{
		Token: token,
		Title: title,
		Data:  buf.Bytes(),
	})
	if err != nil {
		response.PublishErr(c, errno.RpcConnectErr)
	}
	if pubilsh.StatusCode == errno.ServiceErrCode {
		response.PublishErr(c, errno.ServiceErr)
		return
	}
	response.PublishOk(c, errno.Success)
}
