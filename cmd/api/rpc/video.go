package rpc

import (
	"context"
	"fmt"
	"github.com/132982317/profstik/kitex_gen/video"
	"github.com/132982317/profstik/kitex_gen/video/videoservice"
	"github.com/132982317/profstik/pkg/utils/viper"
	"github.com/cloudwego/kitex-examples/bizdemo/easy_note/pkg/middleware"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
	"log"
	"time"
)

var videoClient videoservice.Client

func initVideoRpc(config *viper.Config) {
	// 使用Etcd解析器创建解析器对象
	r, err := etcd.NewEtcdResolver([]string{fmt.Sprintf("%s", config.Viper.GetString("etcd.address"))})
	if err != nil {
		log.Fatal(err)
	}
	serviceName := config.Viper.GetString("server.name")
	// 创建注释rpc客户端
	c, err := videoservice.NewClient(
		serviceName,
		client.WithMiddleware(middleware.CommonMiddleware),
		client.WithMiddleware(middleware.ClientMiddleware),
		client.WithMuxConnection(1),
		client.WithRPCTimeout(3*time.Second),              // rpc超时时间
		client.WithConnectTimeout(50*time.Millisecond),    // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		client.WithSuite(trace.NewDefaultClientSuite()),   // tracer
		client.WithResolver(r),                            // resolver
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serviceName}),
	)
	if err != nil {
		log.Fatal(err)
	}
	videoClient = c
}

func Feed(ctx context.Context, req *video.FeedRequest) (*video.FeedResponse, error) {
	return videoClient.Feed(ctx, req)
}

func PublishAction(ctx context.Context, req *video.PublishActionRequest) (*video.PublishActionResponse, error) {
	return videoClient.PublishAction(ctx, req)
}

func PublishList(ctx context.Context, req *video.PublishListRequest) (*video.PublishListResponse, error) {
	return videoClient.PublishList(ctx, req)
}
