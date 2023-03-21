package mysql

import (
	"context"
	"github.com/132982317/profstik/pkg/errno"
	"gorm.io/gorm"
)

type FavoriteVideoRelation struct {
	Video   Video `gorm:"foreignkey:VideoID;" json:"video,omitempty"`
	VideoID uint  `gorm:"index:idx_videoid;not null" json:"video_id"`
	User    User  `gorm:"foreignkey:UserID;" json:"user,omitempty"`
	UserID  uint  `gorm:"index:idx_userid;not null" json:"user_id"`
}

type FavoriteCommentRelation struct {
	Comment   Comment `gorm:"foreignkey:CommentID;" json:"comment,omitempty"`
	CommentID uint    `gorm:"column:comment_id;index:idx_commentid;not null" json:"video_id"`
	User      User    `gorm:"foreignkey:UserID;" json:"user,omitempty"`
	UserID    uint    `gorm:"column:user_id;index:idx_userid;not null" json:"user_id"`
}

func (FavoriteVideoRelation) TableName() string {
	return "user_favorite_videos"
}

func (FavoriteCommentRelation) TableName() string {
	return "user_favorite_comments"
}

func CreateVideoFavorite(ctx context.Context, userID int64, videoID int64, authorID int64) error {
	err := GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//1. 新增点赞数据
		err := tx.Create(&FavoriteVideoRelation{UserID: uint(userID), VideoID: uint(videoID)}).Error
		if err != nil {
			return err
		}

		//2.改变 video 表中的 favorite count
		res := tx.Model(&Video{}).Where("id = ?", videoID).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			// 影响的数据条数不是1
			return errno.ErrDatabase
		}

		//3.改变 user 表中的 favorite count
		res = tx.Model(&User{}).Where("id = ?", userID).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
		if res.Error != nil {
			return err
		}
		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}

		//4.改变 user 表中的 total_favorited
		res = tx.Model(&User{}).Where("id = ?", authorID).Update("total_favorited", gorm.Expr("total_favorited + ?", 1))
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

func GetFavoriteVideoRelationByUserVideoID(ctx context.Context, userID int64, videoID int64) (*FavoriteVideoRelation, error) {
	FavoriteVideoRelation := new(FavoriteVideoRelation)
	if err := GetDB().WithContext(ctx).First(&FavoriteVideoRelation, "user_id = ? and video_id = ?", userID, videoID).Error; err == nil {
		return FavoriteVideoRelation, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}

func DelFavoriteByUserVideoID(ctx context.Context, userID int64, videoID int64, authorID int64) error {
	err := GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		FavoriteVideoRelation := new(FavoriteVideoRelation)
		if err := tx.Where("user_id = ? and video_id = ?", userID, videoID).First(&FavoriteVideoRelation).Error; err != nil {
			return err
		} else if err == gorm.ErrRecordNotFound {
			return nil
		}

		//1. 删除点赞数据
		err := tx.Unscoped().Where("user_id = ? and video_id = ?", userID, videoID).Delete(&FavoriteVideoRelation).Error
		if err != nil {
			return err
		}

		//2.改变 video 表中的 favorite count
		res := tx.Model(&Video{}).Where("id = ?", videoID).Update("favorite_count", gorm.Expr("favorite_count - ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			// 影响数据条数不是1
			return errno.ErrDatabase
		}

		//3.改变 user 表中的 favorite count
		res = tx.Model(&User{}).Where("id = ?", userID).Update("favorite_count", gorm.Expr("favorite_count - ?", 1))
		if res.Error != nil {
			return err
		}
		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}

		//4.改变 user 表中的 total_favorited
		res = tx.Model(&User{}).Where("id = ?", authorID).Update("total_favorited", gorm.Expr("total_favorited - ?", 1))
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

func GetFavoriteListByUserID(ctx context.Context, userID int64) ([]*FavoriteVideoRelation, error) {
	var FavoriteVideoRelationList []*FavoriteVideoRelation
	err := GetDB().WithContext(ctx).Where("user_id = ?", userID).Find(&FavoriteVideoRelationList).Error
	if err != nil {
		return nil, err
	}
	return FavoriteVideoRelationList, nil
}

func GetAllFavoriteList(ctx context.Context) ([]*FavoriteVideoRelation, error) {
	var FavoriteVideoRelationList []*FavoriteVideoRelation
	if err := GetDB().WithContext(ctx).Find(&FavoriteVideoRelationList).Error; err != nil {
		return nil, err
	}
	return FavoriteVideoRelationList, nil
}

func GetFavoriteCommentRelationByUserCommentID(ctx context.Context, userID int64, commentID int64) (*FavoriteCommentRelation, error) {
	FavoriteCommentRelation := new(FavoriteCommentRelation)
	if err := GetDB().WithContext(ctx).First(&FavoriteCommentRelation, "user_id = ? and comment_id = ?", userID, commentID).Error; err == nil {
		return FavoriteCommentRelation, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}
