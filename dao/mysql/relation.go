package mysql

import (
	"context"
	"github.com/132982317/profstik/pkg/errno"
	"gorm.io/gorm"
)

type FollowRelation struct {
	gorm.Model
	User     User `gorm:"foreignkey:UserID;" json:"user,omitempty"`
	UserID   uint `gorm:"index:idx_userid;not null" json:"user_id"`
	ToUser   User `gorm:"foreignkey:ToUserID;" json:"to_user,omitempty"`
	ToUserID uint `gorm:"index:idx_userid;index:idx_userid_to;not null" json:"to_user_id"`
}

func (FollowRelation) TableName() string {
	return "relations"
}

func GetRelationByUserIDs(ctx context.Context, userID int64, toUserID int64) (*FollowRelation, error) {
	relation := new(FollowRelation)
	if err := GetDB().WithContext(ctx).Where("user_id=? AND to_user_id=?", userID, toUserID).First(&relation).Error; err == nil {
		errno.Dprintf("[GetRelationByUserIDs] relation:%+v", relation)
		return relation, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}

func CreateRelation(ctx context.Context, userID int64, toUserID int64) error {
	err := GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 新增关注数据
		err := tx.Create(&FollowRelation{UserID: uint(userID), ToUserID: uint(toUserID)}).Error
		if err != nil {
			return err
		}
		// 2.改变 user 表中的 following count
		res := tx.Model(&User{}).Where("id = ?", userID).Update("following_count", gorm.Expr("following_count + ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}
		// 3.改变 user 表中的 follower count
		res = tx.Model(&User{}).Where("id = ?", toUserID).Update("follower_count", gorm.Expr("follower_count + ?", 1))
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}
		return nil
	})
	return err
}

func DelRelationByUserIDs(ctx context.Context, userID int64, toUserID int64) error {
	err := GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		relation := new(FollowRelation)
		if err := tx.Where("user_id = ? AND to_user_id=?", userID, toUserID).First(&relation).Error; err != nil {
			return err
		} else if err == gorm.ErrRecordNotFound {
			return nil
		}
		// 1. 删除关注数据
		err := tx.Unscoped().Delete(&relation).Error
		if err != nil {
			return err
		}
		// 2.改变 user 表中的 following count
		res := tx.Model(&User{}).Where("id = ?", userID).Update("following_count", gorm.Expr("following_count - ?", 1))
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}
		// 3.改变 user 表中的 follower count
		res = tx.Model(&User{}).Where("id = ?", toUserID).Update("follower_count", gorm.Expr("follower_count - ?", 1))
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return errno.ErrDatabase
		}
		return nil
	})
	return err
}

// GetFollowingListByUserID
//
//	@Description: 获取指定用户的关注关系列表
//	@Date 2023-01-21 17:01:40
//	@param ctx 数据库操作上下文
//	@param userID 指定用户的用户id
//	@return []*Relation 指定用户的关注关系列表
//	@return error
func GetFollowingListByUserID(ctx context.Context, userID int64) ([]*FollowRelation, error) {
	var RelationList []*FollowRelation
	err := GetDB().WithContext(ctx).Where("user_id = ?", userID).Find(&RelationList).Error
	if err != nil {
		return nil, err
	}
	return RelationList, nil
}

func GetFollowerListByUserID(ctx context.Context, toUserID int64) ([]*FollowRelation, error) {
	var RelationList []*FollowRelation
	err := GetDB().WithContext(ctx).Where("to_user_id = ?", toUserID).Find(&RelationList).Error
	if err != nil {
		return nil, err
	}
	return RelationList, nil
}

func GetFriendList(ctx context.Context, userID int64) ([]*FollowRelation, error) {
	var FriendList []*FollowRelation
	err := GetDB().WithContext(ctx).Raw("SELECT user_id, to_user_id, created_at FROM relations WHERE user_id = ? AND to_user_id IN (SELECT user_id FROM relations r WHERE r.to_user_id = relations.user_id)", userID).Scan(&FriendList).Error
	if err != nil {
		return nil, err
	}
	return FriendList, nil
}
