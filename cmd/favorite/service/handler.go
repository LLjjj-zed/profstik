package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/132982317/profstik/dao/mysql"
	"github.com/132982317/profstik/dao/redis"
	"github.com/132982317/profstik/kitex_gen/favorite"
	"github.com/132982317/profstik/kitex_gen/user"
	"github.com/132982317/profstik/kitex_gen/video"
	"github.com/132982317/profstik/pkg/errno"
	"github.com/132982317/profstik/pkg/utils/minio"
	"github.com/132982317/profstik/pkg/utils/rabbitmq"
	"github.com/132982317/profstik/pkg/utils/zap"
	Zap "go.uber.org/zap"
	"strings"
	"time"
)

type FavoriteServiceImpl struct{}

func (f *FavoriteServiceImpl) FavoriteAction(ctx context.Context, req *favorite.FavoriteActionRequest) (resp *favorite.FavoriteActionResponse, err error) {
	logger := zap.LoggerPool.Get().(*Zap.SugaredLogger)
	defer zap.LoggerPool.Put(logger)
	// 解析token,获取用户id
	claims, err := Jwt.ParseToken(req.Token)
	if err != nil {
		logger.Errorf("token解析错误：%v", err.Error())
		resp = &favorite.FavoriteActionResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	userID := claims.Id

	//将点赞信息存入消息队列,成功存入则表示点赞成功,后续处理由redis完成
	fc := &redis.FavoriteCache{
		VideoID:    uint(req.VideoId),
		UserID:     uint(userID),
		ActionType: uint(req.ActionType),
		CreatedAt:  uint(time.Now().UnixMilli()),
	}
	jsonFC, _ := json.Marshal(fc)
	fmt.Println("Publish new message: ", fc)
	if err = FavoriteMq.PublishSimple(ctx, jsonFC); err != nil {
		logger.Errorf("消息队列发布错误：%v", err.Error())
		if strings.Contains(err.Error(), "连接断开") {
			// 检测到通道关闭，则重连
			go FavoriteMq.Destroy()
			FavoriteMq = rabbitmq.NewRabbitMQSimple("favorite", autoAck)
			logger.Errorln("消息队列通道尝试重连：favorite")
			go consume()
			resp = &favorite.FavoriteActionResponse{
				StatusCode: errno.SuccessCode,
			}
			return resp, nil
		}
		resp = &favorite.FavoriteActionResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	resp = &favorite.FavoriteActionResponse{
		StatusCode: errno.SuccessCode,
	}
	return resp, nil
}

func (f *FavoriteServiceImpl) FavoriteList(ctx context.Context, req *favorite.FavoriteListRequest) (resp *favorite.FavoriteListResponse, err error) {
	logger := zap.LoggerPool.Get().(*Zap.SugaredLogger)
	defer zap.LoggerPool.Put(logger)

	userID := req.UserId
	// 从数据库获取喜欢列表
	results, err := mysql.GetFavoriteListByUserID(ctx, userID)
	if err != nil {
		logger.Errorf("获取喜欢列表错误：%v", err.Error())
		resp = &favorite.FavoriteListResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	favorites := make([]*video.Video, 0)
	for _, r := range results {
		v, err := mysql.GetVideoById(ctx, int64(r.VideoID))
		if err != nil {
			logger.Errorf("获取视频错误：%v", err.Error())
			resp = &favorite.FavoriteListResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}

		u, err := mysql.GetUserByID(ctx, int64(v.AuthorID))
		if err != nil {
			logger.Errorf("获取用户错误：%v", err.Error())
			resp = &favorite.FavoriteListResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}

		relation, err := mysql.GetRelationByUserIDs(ctx, userID, int64(u.ID))
		if err != nil {
			logger.Errorf("发生错误：%v", err.Error())
			resp = &favorite.FavoriteListResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
		playUrl, err := minio.GetFileTemporaryURL(minio.VideoBucketName, v.PlayUrl)
		if err != nil {
			logger.Errorf("发生错误：%v", err.Error())
			resp = &favorite.FavoriteListResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
		coverUrl, err := minio.GetFileTemporaryURL(minio.CoverBucketName, v.CoverUrl)
		if err != nil {
			logger.Errorf("发生错误：%v", err.Error())
			resp = &favorite.FavoriteListResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
		avatar, err := minio.GetFileTemporaryURL(minio.AvatarBucketName, u.Avatar)
		if err != nil {
			logger.Errorf("Minio获取头像失败：%v", err.Error())
			resp = &favorite.FavoriteListResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
		backgroundUrl, err := minio.GetFileTemporaryURL(minio.BackgroundImageBucketName, u.BackgroundImage)
		if err != nil {
			logger.Errorf("Minio获取链接失败：%v", err.Error())
			resp = &favorite.FavoriteListResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
		favorites = append(favorites, &video.Video{
			Id: int64(r.VideoID),
			Author: &user.User{
				Id:              int64(u.ID),
				Name:            u.UserName,
				FollowCount:     int64(u.FollowingCount),
				FollowerCount:   int64(u.FollowerCount),
				IsFollow:        relation != nil,
				Avatar:          avatar,
				BackgroundImage: backgroundUrl,
				Signature:       u.Signature,
				TotalFavorited:  int64(u.TotalFavorited),
				WorkCount:       int64(u.WorkCount),
				FavoriteCount:   int64(u.FavoriteCount),
			},
			PlayUrl:       playUrl,
			CoverUrl:      coverUrl,
			FavoriteCount: int64(v.FavoriteCount),
			CommentCount:  int64(v.CommentCount),
			IsFavorite:    true,
			Title:         v.Title,
		})
	}

	resp = &favorite.FavoriteListResponse{
		StatusCode: errno.SuccessCode,
		VideoList:  favorites,
	}
	return resp, nil
}
