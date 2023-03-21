package mysql

import (
	"context"
	"github.com/132982317/profstik/pkg/errno"
	"gorm.io/gorm"
	"time"
)

type Video struct {
	ID            uint      `gorm:"primarykey"`
	CreatedAt     time.Time `gorm:"not null;index:idx_create" json:"created_at,omitempty"`
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Author        User           `gorm:"foreignkey:AuthorID" json:"author,omitempty"`
	AuthorID      uint           `gorm:"index:idx_authorid;not null" json:"author_id,omitempty"`
	PlayUrl       string         `gorm:"type:varchar(255);not null" json:"play_url,omitempty"`
	CoverUrl      string         `gorm:"type:varchar(255)" json:"cover_url,omitempty"`
	FavoriteCount uint           `gorm:"default:0;not null" json:"favorite_count,omitempty"`
	CommentCount  uint           `gorm:"default:0;not null" json:"comment_count,omitempty"`
	Title         string         `gorm:"type:varchar(50);not null" json:"title,omitempty"`
}

func (Video) TableName() string {
	return "videos"
}

func MGetVideos(ctx context.Context, limit int, latestTime *int64) ([]*Video, error) {
	videos := make([]*Video, 0)

	if latestTime == nil || *latestTime == 0 {
		curTime := time.Now().UnixMilli()
		latestTime = &curTime
	}
	if err := GetDB().WithContext(ctx).Limit(limit).Order("created_at desc").Find(&videos, "created_at < ?", time.UnixMilli(*latestTime)).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

func GetVideoById(ctx context.Context, videoID int64) (*Video, error) {
	video := new(Video)
	if err := GetDB().WithContext(ctx).First(&video, videoID).Error; err == nil {
		errno.Dprintf("[GetVideoById] video:%+v", video)
		return video, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}

func GetVideoListByIDs(ctx context.Context, videoIDs []int64) ([]*Video, error) {
	videos := make([]*Video, 0)
	if len(videoIDs) == 0 {
		return videos, nil
	}
	if err := GetDB().WithContext(ctx).Where("video_id in ?", videoIDs).Find(&videos).Error; err != nil {
		return nil, err
	}
	errno.Dprintf("[GetVideoById] videos:%+v", videos)
	return videos, nil
}

func CreateVideo(ctx context.Context, video *Video) error {
	err := GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 在 video 表中创建视频记录
		err := tx.Create(video).Error
		if err != nil {
			return err
		}
		// 2. 同步 user 表中的作品数量
		res := tx.Model(&User{}).Where("id = ?", video.AuthorID).Update("work_count", gorm.Expr("work_count + ?", 1))
		if res.Error != nil {
			return err
		}
		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}
		return nil
	})

	return err
}

func GetVideosByUserID(ctx context.Context, authorId int64) ([]*Video, error) {
	var pubList []*Video
	err := GetDB().WithContext(ctx).Model(&Video{}).Where(&Video{AuthorID: uint(authorId)}).Find(&pubList).Error
	if err != nil {
		return nil, err
	}
	return pubList, nil
}

func DelVideoByID(ctx context.Context, videoID int64, authorID int64) error {
	err := GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 根据主键 video_id 删除 video
		err := tx.Unscoped().Delete(&Video{}, videoID).Error
		if err != nil {
			return err
		}
		// 2. 同步 user 表中的作品数量
		res := tx.Model(&User{}).Where("id = ?", authorID).Update("work_count", gorm.Expr("work_count - ?", 1))
		if res.Error != nil {
			return err
		}
		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}
		return nil
	})
	return err
}
