package rpc

import (
	"context"
	"fmt"
	"github.com/132982317/profstik/kitex_gen/relation"
	"github.com/132982317/profstik/kitex_gen/relation/relationservice"
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

var relationClient relationservice.Client

func initRelationRpc(config *viper.Config) {
	// 使用Etcd解析器创建解析器对象
	r, err := etcd.NewEtcdResolver([]string{fmt.Sprintf("%s", config.Viper.GetString("etcd.address"))})
	if err != nil {
		log.Fatal(err)
	}
	serviceName := config.Viper.GetString("server.name")
	// 创建注释rpc客户端
	c, err := relationservice.NewClient(
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
	relationClient = c
}

func RelationAction(ctx context.Context, req *relation.RelationActionRequest) (*relation.RelationActionResponse, error) {
	return relationClient.RelationAction(ctx, req)
}

func RelationFollowList(ctx context.Context, req *relation.RelationFollowListRequest) (*relation.RelationFollowListResponse, error) {
	return relationClient.RelationFollowList(ctx, req)
}

func RelationFollowerList(ctx context.Context, req *relation.RelationFollowerListRequest) (*relation.RelationFollowerListResponse, error) {
	return relationClient.RelationFollowerList(ctx, req)
}

func RelationFriendList(ctx context.Context, req *relation.RelationFriendListRequest) (*relation.RelationFriendListResponse, error) {
	return relationClient.RelationFriendList(ctx, req)
}
