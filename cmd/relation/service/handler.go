package service

import (
	"context"
	"encoding/json"
	"github.com/132982317/profstik/dao/mysql"
	"github.com/132982317/profstik/dao/redis"
	"github.com/132982317/profstik/kitex_gen/relation"
	"github.com/132982317/profstik/kitex_gen/user"
	"github.com/132982317/profstik/pkg/errno"
	tool "github.com/132982317/profstik/pkg/utils/crypt"
	"github.com/132982317/profstik/pkg/utils/minio"
	"github.com/132982317/profstik/pkg/utils/rabbitmq"
	"github.com/132982317/profstik/pkg/utils/zap"
	Zap "go.uber.org/zap"
	"strings"
	"time"
)

type RelationServiceImpl struct{}

func (r *RelationServiceImpl) RelationAction(ctx context.Context, req *relation.RelationActionRequest) (resp *relation.RelationActionResponse, err error) {
	logger := zap.LoggerPool.Get().(*Zap.SugaredLogger)
	defer zap.LoggerPool.Put(logger)
	// 解析token,获取用户id
	claims, err := Jwt.ParseToken(req.Token)
	if err != nil {
		logger.Errorln(err.Error())
		resp = &relation.RelationActionResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	userID := claims.Id
	toUserID := req.ToUserId

	if userID == toUserID {
		logger.Errorf("操作非法：用户无法成为自己的粉丝：%d", userID)
		resp = &relation.RelationActionResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	if req.ActionType != 1 && req.ActionType != 2 {
		logger.Errorln("action_type 格式错误")
		resp = &relation.RelationActionResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	// 检查ID是否存在
	u1, _ := mysql.GetUserByID(ctx, userID)
	u2, _ := mysql.GetUserByID(ctx, toUserID)
	if u1 == nil || u2 == nil {
		logger.Errorln("所请求的用户ID不存在")
		resp = &relation.RelationActionResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}

	// 将关注信息存入消息队列，成功存入则表示操作成功，后续处理由redis完成
	relationCache := &redis.RelationCache{
		UserID:     uint(userID),
		ToUserID:   uint(toUserID),
		ActionType: uint(req.ActionType),
		CreatedAt:  uint(time.Now().UnixMilli()),
	}
	jsonRc, _ := json.Marshal(relationCache)
	if err = RelationMq.PublishSimple(ctx, jsonRc); err != nil {
		logger.Errorf("消息队列发布错误：%v", err.Error())
		if strings.Contains(err.Error(), "连接断开") {
			// 检测到通道关闭，则重连
			go RelationMq.Destroy()
			RelationMq = rabbitmq.NewRabbitMQSimple("relation", autoAck)
			logger.Errorln("消息队列通道尝试重连：relation")
			go consume()
			resp = &relation.RelationActionResponse{
				StatusCode: errno.SuccessCode,
			}
			return resp, nil
		}
		resp = &relation.RelationActionResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	resp = &relation.RelationActionResponse{
		StatusCode: errno.SuccessCode,
	}
	return resp, nil
}

func (r *RelationServiceImpl) RelationFollowList(ctx context.Context, req *relation.RelationFollowListRequest) (resp *relation.RelationFollowListResponse, err error) {
	logger := zap.LoggerPool.Get().(*Zap.SugaredLogger)
	defer zap.LoggerPool.Put(logger)

	userID := req.UserId
	// 解析token,获取用户id
	claims, err := Jwt.ParseToken(req.Token)
	if err != nil {
		logger.Errorln(err.Error())
		resp = &relation.RelationFollowListResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	if userID != claims.Id {
		logger.Errorf("当前登录用户%d无法访问其他用户的关注列表%d", claims.Id, userID)
		resp = &relation.RelationFollowListResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}

	// 从数据库获取关注列表
	followings, err := mysql.GetFollowingListByUserID(ctx, userID)
	if err != nil {
		logger.Errorln(err.Error())
		resp = &relation.RelationFollowListResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	userIDs := make([]int64, 0)
	for _, res := range followings {
		userIDs = append(userIDs, int64(res.ToUserID))
	}
	users, err := mysql.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		logger.Errorln(err.Error())
		resp = &relation.RelationFollowListResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	userList := make([]*user.User, 0)
	for _, u := range users {
		avatar, err := minio.GetFileTemporaryURL(minio.AvatarBucketName, u.Avatar)
		if err != nil {
			logger.Errorf("Minio获取头像失败：%v", err.Error())
			resp = &relation.RelationFollowListResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
		backgroundUrl, err := minio.GetFileTemporaryURL(minio.BackgroundImageBucketName, u.BackgroundImage)
		if err != nil {
			logger.Errorf("Minio获取背景图链接失败：%v", err.Error())
			resp = &relation.RelationFollowListResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
		userList = append(userList, &user.User{
			Id:              int64(u.ID),
			Name:            u.UserName,
			FollowCount:     int64(u.FollowingCount),
			FollowerCount:   int64(u.FollowerCount),
			IsFollow:        true,
			Avatar:          avatar,
			BackgroundImage: backgroundUrl,
			Signature:       u.Signature,
			TotalFavorited:  int64(u.TotalFavorited),
			WorkCount:       int64(u.WorkCount),
			FavoriteCount:   int64(u.FavoriteCount),
		})
	}

	// 返回结果
	resp = &relation.RelationFollowListResponse{
		StatusCode: errno.SuccessCode,
		UserList:   userList,
	}
	return resp, nil
}

func (r *RelationServiceImpl) RelationFollowerList(ctx context.Context, req *relation.RelationFollowerListRequest) (resp *relation.RelationFollowerListResponse, err error) {
	logger := zap.LoggerPool.Get().(*Zap.SugaredLogger)
	defer zap.LoggerPool.Put(logger)

	userID := req.UserId

	// 解析token,获取用户id
	claims, err := Jwt.ParseToken(req.Token)
	if err != nil {
		logger.Errorln(err.Error())
		resp = &relation.RelationFollowerListResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	if userID != claims.Id {
		logger.Errorf("当前登录用户%d无法访问其他用户的粉丝列表%d", claims.Id, userID)
		resp = &relation.RelationFollowerListResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}

	// 从数据库获取粉丝列表
	followers, err := mysql.GetFollowerListByUserID(ctx, userID)
	if err != nil {
		logger.Errorln(err.Error())
		resp = &relation.RelationFollowerListResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	userIDs := make([]int64, 0)
	for _, res := range followers {
		userIDs = append(userIDs, int64(res.UserID))
	}
	users, err := mysql.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		logger.Errorln(err.Error())
		resp = &relation.RelationFollowerListResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	userList := make([]*user.User, 0)
	for _, u := range users {
		// 查询两个用户是否互相关注
		follow, err := mysql.GetRelationByUserIDs(ctx, userID, int64(u.ID))
		if err != nil {
			logger.Errorln(err.Error())
			res := &relation.RelationFollowerListResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return res, nil
		}
		avatar, err := minio.GetFileTemporaryURL(minio.AvatarBucketName, u.Avatar)
		if err != nil {
			logger.Errorf("Minio获取头像失败：%v", err.Error())
			resp = &relation.RelationFollowerListResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
		backgroundUrl, err := minio.GetFileTemporaryURL(minio.BackgroundImageBucketName, u.BackgroundImage)
		if err != nil {
			logger.Errorf("Minio获取背景图链接失败：%v", err.Error())
			resp = &relation.RelationFollowerListResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
		userList = append(userList, &user.User{
			Id:              int64(u.ID),
			Name:            u.UserName,
			FollowCount:     int64(u.FollowingCount),
			FollowerCount:   int64(u.FollowerCount),
			IsFollow:        follow != nil,
			Avatar:          avatar,
			BackgroundImage: backgroundUrl,
			Signature:       u.Signature,
			TotalFavorited:  int64(u.TotalFavorited),
			WorkCount:       int64(u.WorkCount),
			FavoriteCount:   int64(u.FavoriteCount),
		})
	}

	// 返回结果
	resp = &relation.RelationFollowerListResponse{
		StatusCode: errno.SuccessCode,
		UserList:   userList,
	}
	return resp, nil
}

func (r *RelationServiceImpl) RelationFriendList(ctx context.Context, req *relation.RelationFriendListRequest) (resp *relation.RelationFriendListResponse, err error) {
	logger := zap.LoggerPool.Get().(*Zap.SugaredLogger)
	defer zap.LoggerPool.Put(logger)
	userID := req.UserId

	// 解析token,获取用户id
	claims, err := Jwt.ParseToken(req.Token)
	if err != nil {
		logger.Errorln(err.Error())
		resp = &relation.RelationFriendListResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	if userID != claims.Id {
		logger.Errorf("当前登录用户%d无法访问其他用户的朋友列表%d", claims.Id, userID)
		resp = &relation.RelationFriendListResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}

	// 从数据库获取朋友列表
	friends, err := mysql.GetFriendList(ctx, userID)
	if err != nil {
		logger.Errorln(err.Error())
		resp = &relation.RelationFriendListResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	userIDs := make([]int64, 0)
	for _, res := range friends {
		userIDs = append(userIDs, int64(res.ToUserID))
	}
	users, err := mysql.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		logger.Errorln(err.Error())
		resp = &relation.RelationFriendListResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	userList := make([]*relation.FriendUser, 0)
	for _, u := range users {
		message, err := mysql.GetFriendLatestMessage(ctx, userID, int64(u.ID))
		if err != nil {
			resp = &relation.RelationFriendListResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
		var msgType int64
		if int64(message.FromUserID) == userID {
			// 当前用户为发送方
			msgType = 1
		} else {
			// 当前用户为接收方
			msgType = 0
		}
		var decContent []byte
		if len(message.Content) != 0 {
			decContent, err = tool.Base64Decode([]byte(message.Content))
			if err != nil {
				logger.Errorf("Base64Decode error: %v\n", err.Error())
				resp = &relation.RelationFriendListResponse{
					StatusCode: errno.ServiceErrCode,
				}
				return
			}
			decContent, err = tool.RsaDecrypt(decContent, privateKey)
			if err != nil {
				logger.Errorf("rsa decrypt error: %v\n", err.Error())
				resp = &relation.RelationFriendListResponse{
					StatusCode: errno.ServiceErrCode,
				}
				return
			}
		}
		avatar, err := minio.GetFileTemporaryURL(minio.AvatarBucketName, u.Avatar)
		if err != nil {
			logger.Errorf("Minio获取头像失败：%v", err.Error())
			resp = &relation.RelationFriendListResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
		backgroundUrl, err := minio.GetFileTemporaryURL(minio.BackgroundImageBucketName, u.BackgroundImage)
		if err != nil {
			logger.Errorf("Minio获取背景图失败：%v", err.Error())
			resp = &relation.RelationFriendListResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
		userList = append(userList, &relation.FriendUser{
			Id:              int64(u.ID),
			Name:            u.UserName,
			FollowCount:     int64(u.FollowingCount),
			FollowerCount:   int64(u.FollowerCount),
			IsFollow:        true,
			Message:         string(decContent),
			MsgType:         msgType,
			Avatar:          avatar,
			BackgroundImage: backgroundUrl,
			Signature:       u.Signature,
			TotalFavorited:  int64(u.TotalFavorited),
			WorkCount:       int64(u.WorkCount),
			FavoriteCount:   int64(u.FavoriteCount),
		})
	}

	// 返回结果
	resp = &relation.RelationFriendListResponse{
		StatusCode: errno.SuccessCode,
		UserList:   userList,
	}
	return resp, nil
}
