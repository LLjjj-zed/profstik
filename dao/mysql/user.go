package mysql

import (
	"context"
	"github.com/132982317/profstik/pkg/errno"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName        string  `gorm:"index:idx_username,unique;type:varchar(40);not null" json:"name,omitempty"`
	Password        string  `gorm:"type:varchar(256);not null" json:"password,omitempty"`
	FavoriteVideos  []Video `gorm:"many2many:user_favorite_videos" json:"favorite_videos,omitempty"`
	FollowingCount  uint    `gorm:"default:0;not null" json:"follow_count,omitempty"`                                                           // 关注总数
	FollowerCount   uint    `gorm:"default:0;not null" json:"follower_count,omitempty"`                                                         // 粉丝总数
	Avatar          string  `gorm:"type:varchar(256)" json:"avatar,omitempty"`                                                                  // 用户头像
	BackgroundImage string  `gorm:"column:background_image;type:varchar(256);default:default_background.jpg" json:"background_image,omitempty"` // 用户个人页顶部大图
	WorkCount       uint    `gorm:"default:0;not null" json:"work_count,omitempty"`                                                             // 作品数
	FavoriteCount   uint    `gorm:"default:0;not null" json:"favorite_count,omitempty"`                                                         // 喜欢数
	TotalFavorited  uint    `gorm:"default:0;not null" json:"total_favorited,omitempty"`                                                        // 获赞总量
	Signature       string  `gorm:"type:varchar(256)" json:"signature,omitempty"`                                                               // 个人简介
}

func (User) TableName() string {
	return "users"
}

func CreateUser(ctx context.Context, user *User) error {
	err := GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		return nil
	})
	errno.Dprintf("[CreateUser] user:%+v", user)
	return err
}

func GetUserByUserName(ctx context.Context, username string) (*User, error) {
	user := new(User)
	if err := GetDB().WithContext(ctx).Select("id, user_name, password").Where("user_name=?", username).First(&user).Error; err == nil {
		errno.Dprintf("[GetUserByUserName] user:%+v", user)
		return user, nil
	} else if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else {
		return nil, err
	}
}

func GetUserByID(ctx context.Context, userid int64) (*User, error) {
	user := new(User)
	if err := GetDB().WithContext(ctx).Where("user_id=?", userid).First(&user).Error; err != nil {
		errno.Dprintf("[GetUserByID] user:%+v", user)
		return nil, err
	}
	return user, nil
}

func GetUsersByIDs(ctx context.Context, userIDs []int64) ([]*User, error) {
	users := make([]*User, 0)
	if len(userIDs) == 0 {
		return users, nil
	}
	if err := GetDB().WithContext(ctx).Where("id in ?", userIDs).Find(&users).Error; err != nil {
		errno.Dprintf("[GetUserByIDs] users:%+v", users)
		return nil, err
	}
	return users, nil
}
