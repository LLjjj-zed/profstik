package service

import (
	"context"
	"fmt"
	"github.com/132982317/profstik/dao/mysql"
	"github.com/132982317/profstik/kitex_gen/user"
	"github.com/132982317/profstik/kitex_gen/video"
	"github.com/132982317/profstik/pkg/errno"
	"github.com/132982317/profstik/pkg/utils/minio"
	"github.com/132982317/profstik/pkg/utils/viper"
	"github.com/132982317/profstik/pkg/utils/zap"
	Zap "go.uber.org/zap"
	"time"
)

type VideoServiceImpl struct{}

const limit = 30

func (v *VideoServiceImpl) Feed(ctx context.Context, req *video.FeedRequest) (resp *video.FeedResponse, err error) {
	logger := zap.LoggerPool.Get().(*Zap.SugaredLogger)
	defer zap.LoggerPool.Put(logger)

	nextTime := time.Now().UnixMilli()
	var userID int64 = -1

	// 验证token有效性
	if req.Token != "" {
		claims, err := Jwt.ParseToken(req.Token)
		if err != nil {
			logger.Errorln(err.Error())
			resp = &video.FeedResponse{
				StatusCode: errno.AuthorizationFailedErrCode,
			}
			return
		}
		userID = claims.Id
	}
	// 调用数据库查询 video_list
	videos, err := mysql.MGetVideos(ctx, limit, &req.LatestTime)
	if err != nil {
		logger.Errorln(err.Error())
		resp = &video.FeedResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	videoList := make([]*video.Video, 0)
	for _, r := range videos {
		author, err := mysql.GetUserByID(ctx, int64(r.AuthorID))
		if err != nil {
			logger.Errorf("error:%v", err.Error())
			return nil, err
		}
		relation, err := mysql.GetRelationByUserIDs(ctx, userID, int64(author.ID))
		if err != nil {
			logger.Errorln(err.Error())
			resp = &video.FeedResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
		favorite, err := mysql.GetFavoriteVideoRelationByUserVideoID(ctx, userID, int64(r.ID))
		if err != nil {
			logger.Errorln(err.Error())
			resp = &video.FeedResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
		playUrl, err := minio.GetFileTemporaryURL(minio.VideoBucketName, r.PlayUrl)
		if err != nil {
			logger.Errorf("Minio获取链接失败：%v", err.Error())
			resp = &video.FeedResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
		coverUrl, err := minio.GetFileTemporaryURL(minio.CoverBucketName, r.CoverUrl)
		if err != nil {
			logger.Errorf("Minio获取链接失败：%v", err.Error())
			resp = &video.FeedResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
		avatarUrl, err := minio.GetFileTemporaryURL(minio.AvatarBucketName, author.Avatar)
		if err != nil {
			logger.Errorf("Minio获取链接失败：%v", err.Error())
			resp = &video.FeedResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
		backgroundUrl, err := minio.GetFileTemporaryURL(minio.BackgroundImageBucketName, author.BackgroundImage)
		if err != nil {
			logger.Errorf("Minio获取链接失败：%v", err.Error())
			resp = &video.FeedResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}

		videoList = append(videoList, &video.Video{
			Id: int64(r.ID),
			Author: &user.User{
				Id:              int64(author.ID),
				Name:            author.UserName,
				FollowCount:     int64(author.FollowingCount),
				FollowerCount:   int64(author.FollowerCount),
				IsFollow:        relation != nil,
				Avatar:          avatarUrl,
				BackgroundImage: backgroundUrl,
				Signature:       author.Signature,
				TotalFavorited:  int64(author.TotalFavorited),
				WorkCount:       int64(author.WorkCount),
				FavoriteCount:   int64(author.FavoriteCount),
			},
			PlayUrl:       playUrl,
			CoverUrl:      coverUrl,
			FavoriteCount: int64(r.FavoriteCount),
			CommentCount:  int64(r.CommentCount),
			IsFavorite:    favorite != nil,
			Title:         r.Title,
		})
	}
	if len(videos) != 0 {
		nextTime = videos[len(videos)-1].UpdatedAt.UnixMilli()
	}
	resp = &video.FeedResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		VideoList:  videoList,
		NextTime:   nextTime,
	}
	return resp, nil
}

func (v *VideoServiceImpl) PublishAction(ctx context.Context, req *video.PublishActionRequest) (resp *video.PublishActionResponse, err error) {
	logger := zap.LoggerPool.Get().(*Zap.SugaredLogger)
	defer zap.LoggerPool.Put(logger)
	// 解析token,获取用户id
	claims, err := Jwt.ParseToken(req.Token)
	if err != nil {
		logger.Errorln(err.Error())
		resp = &video.PublishActionResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	userID := claims.Id

	if len(req.Title) == 0 || len(req.Title) > 32 {
		logger.Errorf("标题不能为空且不能超过32个字符：%d", len(req.Title))
		resp = &video.PublishActionResponse{
			StatusCode: errno.ServiceErrCode,
			StatusMsg:  "标题不能为空且不能超过32个字符",
		}
		return
	}

	// 限制文件上传大小
	maxSize := viper.Init("video").Viper.GetInt("video.maxSizeLimit")
	size := len(req.Data)
	if size > maxSize*1000*1000 {
		logger.Errorln("视频文件过大")
		resp = &video.PublishActionResponse{
			StatusCode: errno.ServiceErrCode,
			StatusMsg:  fmt.Sprintf("该视频文件大于%dMB，上传受限", maxSize),
		}
		return
	}

	createTimestamp := time.Now().UnixMilli()
	videoTitle, coverTitle := fmt.Sprintf("%d_%s_%d.mp4", userID, req.Title, createTimestamp), fmt.Sprintf("%d_%s_%d.png", userID, req.Title, createTimestamp)

	// 插入数据库
	videos := &mysql.Video{
		Title:    req.Title,
		PlayUrl:  videoTitle,
		CoverUrl: coverTitle,
		AuthorID: uint(userID),
	}
	err = mysql.CreateVideo(ctx, videos)
	if err != nil {
		logger.Errorln(err.Error())
		resp = &video.PublishActionResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}

	go func() {
		err := VideoPublish(req.Data, videoTitle, coverTitle)
		if err != nil {
			// 发生错误，则删除插入的记录
			e := mysql.DelVideoByID(ctx, int64(videos.ID), userID)
			if e != nil {
				logger.Errorf("视频记录删除失败：%s", err.Error())
			}
		}
	}()

	resp = &video.PublishActionResponse{
		StatusCode: errno.ServiceErrCode,
	}
	return resp, nil
}

func (v *VideoServiceImpl) PublishList(ctx context.Context, req *video.PublishListRequest) (resp *video.PublishListResponse, err error) {
	logger := zap.LoggerPool.Get().(*Zap.SugaredLogger)
	defer zap.LoggerPool.Put(logger)

	userID := req.UserId
	results, err := mysql.GetVideosByUserID(ctx, userID)
	if err != nil {
		logger.Errorln(err.Error())
		resp = &video.PublishListResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	videos := make([]*video.Video, 0)
	for _, r := range results {
		author, err := mysql.GetUserByID(ctx, int64(r.AuthorID))
		if err != nil {
			logger.Errorln(err.Error())
			resp = &video.PublishListResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
		follow, err := mysql.GetRelationByUserIDs(ctx, userID, int64(author.ID))
		if err != nil {
			logger.Errorln(err.Error())
			resp = &video.PublishListResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
		favorite, err := mysql.GetFavoriteVideoRelationByUserVideoID(ctx, userID, int64(r.ID))
		if err != nil {
			logger.Errorln(err.Error())
			resp = &video.PublishListResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
		playUrl, err := minio.GetFileTemporaryURL(minio.VideoBucketName, r.PlayUrl)
		if err != nil {
			logger.Errorln(err.Error())
			resp = &video.PublishListResponse{
				StatusCode: errno.SuccessCode,
			}
			return
		}
		coverUrl, err := minio.GetFileTemporaryURL(minio.CoverBucketName, r.CoverUrl)
		if err != nil {
			logger.Errorln(err.Error())
			resp = &video.PublishListResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
		avatarUrl, err := minio.GetFileTemporaryURL(minio.AvatarBucketName, author.Avatar)
		if err != nil {
			logger.Errorln(err.Error())
			resp = &video.PublishListResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
		backgroundUrl, err := minio.GetFileTemporaryURL(minio.BackgroundImageBucketName, author.BackgroundImage)
		if err != nil {
			logger.Errorln(err.Error())
			resp = &video.PublishListResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}

		videos = append(videos, &video.Video{
			Id: int64(r.ID),
			Author: &user.User{
				Id:              int64(author.ID),
				Name:            author.UserName,
				FollowerCount:   int64(author.FollowerCount),
				FollowCount:     int64(author.FollowingCount),
				IsFollow:        follow != nil,
				Avatar:          avatarUrl,
				BackgroundImage: backgroundUrl,
				Signature:       author.Signature,
				TotalFavorited:  int64(author.TotalFavorited),
				WorkCount:       int64(author.WorkCount),
				FavoriteCount:   int64(author.FavoriteCount),
			},
			PlayUrl:       playUrl,
			CoverUrl:      coverUrl,
			FavoriteCount: int64(r.FavoriteCount),
			CommentCount:  int64(r.CommentCount),
			IsFavorite:    favorite != nil,
			Title:         r.Title,
		})
	}

	resp = &video.PublishListResponse{
		StatusCode: errno.ServiceErrCode,
		VideoList:  videos,
	}
	return resp, nil
}
