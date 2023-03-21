package service

import (
	"context"
	"github.com/132982317/profstik/dao/mysql"
	"github.com/132982317/profstik/kitex_gen/comment"
	"github.com/132982317/profstik/kitex_gen/user"
	"github.com/132982317/profstik/pkg/errno"
	"github.com/132982317/profstik/pkg/utils/zap"
	Zap "go.uber.org/zap"
	"gorm.io/gorm"
)

type CommentServiceImpl struct{}

func (c *CommentServiceImpl) CommentAction(ctx context.Context, req *comment.CommentActionRequest) (resp *comment.CommentActionResponse, err error) {
	logger := zap.LoggerPool.Get().(*Zap.SugaredLogger)
	defer zap.LoggerPool.Put(logger)
	// 解析token,获取用户id
	claims, err := Jwt.ParseToken(req.Token)
	if err != nil {
		logger.Errorf("token解析错误：%v", err.Error())
		resp = &comment.CommentActionResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	userID := claims.Id
	actionType := req.ActionType
	v, _ := mysql.GetVideoById(ctx, req.VideoId)
	if v == nil {
		logger.Errorf("该视频ID不存在：%d", req.VideoId)
		resp = &comment.CommentActionResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	if actionType == 1 {
		cmt := &mysql.Comment{
			VideoID: uint(req.VideoId),
			UserID:  uint(userID),
			Content: req.CommentText,
		}
		err := mysql.CreateComment(ctx, cmt)
		if err != nil {
			logger.Errorf("新增评论失败：%v", err.Error())
			resp = &comment.CommentActionResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
	} else if actionType == 2 {
		// 判断该评论是否发布自该用户，或该评论在该用户所发布的视频下
		cmt, err := mysql.GetCommentByCommentID(ctx, req.CommentId)
		if err != nil {
			logger.Errorf("评论删除失败：%v", err.Error())
			res := &comment.CommentActionResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return res, nil
		}
		if cmt == nil {
			// 评论不存在，无法删除
			logger.Errorf("评论删除失败，该评论ID不存在：%v", req.CommentId)
			resp = &comment.CommentActionResponse{
				StatusCode: errno.ServiceErrCode,
			}
		} else {
			// 查找该视频的作者ID
			v, err := mysql.GetVideoById(ctx, int64(cmt.VideoID))
			if err != nil {
				logger.Errorf("评论删除失败：%v", err.Error())
				resp = &comment.CommentActionResponse{
					StatusCode: errno.ServiceErrCode,
				}
				return
			}
			// 若删除评论的用户不是发布评论的用户或该用户不是视频创作者
			if userID != int64(cmt.UserID) || userID != int64(v.AuthorID) {
				logger.Errorf("评论删除失败，没有权限：%v", cmt.UserID)
				resp = &comment.CommentActionResponse{
					StatusCode: errno.ServiceErrCode,
				}
				return
			}
		}
		err = mysql.DelCommentByID(ctx, req.CommentId, req.VideoId)
		if err != nil {
			logger.Errorf("评论删除失败：%v", err.Error())
			resp = &comment.CommentActionResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
	} else {
		resp = &comment.CommentActionResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	resp = &comment.CommentActionResponse{
		StatusCode: errno.SuccessCode,
	}
	return resp, nil
}

func (c *CommentServiceImpl) CommentList(ctx context.Context, req *comment.CommentListRequest) (resp *comment.CommentListResponse, err error) {
	logger := zap.LoggerPool.Get().(*Zap.SugaredLogger)
	defer zap.LoggerPool.Put(logger)
	var userID int64 = -1
	// 验证token有效性
	if req.Token != "" {
		claims, err := Jwt.ParseToken(req.Token)
		if err != nil {
			logger.Errorf("token解析错误:%v", err)
			resp = &comment.CommentListResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
		userID = claims.Id
	}

	// 从数据库获取评论列表
	results, err := mysql.GetVideoCommentListByVideoID(ctx, req.VideoId)
	if err != nil {
		logger.Errorf("获取评论列表错误：%v", err)
		resp = &comment.CommentListResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	comments := make([]*comment.Comment, 0)
	for _, r := range results {
		u, err := mysql.GetUserByID(ctx, int64(r.UserID))
		if err != nil && err != gorm.ErrRecordNotFound {
			logger.Errorf("获取用户错误：%v", err.Error())
			resp = &comment.CommentListResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
		_, err = mysql.GetRelationByUserIDs(ctx, userID, int64(u.ID))
		if err != nil {
			logger.Errorln(err.Error())
			resp = &comment.CommentListResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
		usr := &user.User{
			Id:              userID,
			Name:            u.UserName,
			FollowCount:     int64(u.FollowingCount),
			FollowerCount:   int64(u.FollowerCount),
			IsFollow:        err != gorm.ErrRecordNotFound,
			Avatar:          "avatar",
			BackgroundImage: "backgroundUrl",
			Signature:       u.Signature,
			TotalFavorited:  int64(u.TotalFavorited),
			WorkCount:       int64(u.WorkCount),
			FavoriteCount:   int64(u.FavoriteCount),
		}
		comments = append(comments, &comment.Comment{
			Id:         int64(r.ID),
			User:       usr,
			Content:    r.Content,
			CreateDate: r.CreatedAt.Format("2006-01-02"),
			LikeCount:  int64(r.LikeCount),
			TeaseCount: int64(r.TeaseCount),
		})
	}

	resp = &comment.CommentListResponse{
		StatusCode:  errno.SuccessCode,
		CommentList: comments,
	}
	return resp, nil
}
