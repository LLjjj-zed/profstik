package mysql

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type Message struct {
	ID         uint      `gorm:"primarykey"`
	CreatedAt  time.Time `gorm:"index;not null" json:"create_time"`
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	FromUser   User           `gorm:"foreignkey:FromUserID;" json:"from_user,omitempty"`
	FromUserID uint           `gorm:"index:idx_userid_from;not null" json:"from_user_id"`
	ToUser     User           `gorm:"foreignkey:ToUserID;" json:"to_user,omitempty"`
	ToUserID   uint           `gorm:"index:idx_userid_from;index:idx_userid_to;not null" json:"to_user_id"`
	Content    string         `gorm:"type:varchar(255);not null" json:"content"`
}

func (Message) TableName() string {
	return "messages"
}

func GetMessagesByUserIDs(ctx context.Context, userID int64, toUserID int64, lastTimestamp int64) ([]*Message, error) {
	res := make([]*Message, 0)
	if err := GetDB().WithContext(ctx).Where("((from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)) AND created_at > ?",
		userID, toUserID, toUserID, userID, time.UnixMilli(lastTimestamp).Format("2006-01-02 15:04:05.000"),
	).Order("created_at ASC").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func GetMessagesByUserToUser(ctx context.Context, userID int64, toUserID int64, lastTimestamp int64) ([]*Message, error) {
	res := make([]*Message, 0)
	if err := GetDB().WithContext(ctx).Where("from_user_id = ? AND to_user_id = ? AND created_at > ?",
		userID, toUserID, time.UnixMilli(lastTimestamp).Format("2006-01-02 15:04:05.000"),
	).Order("created_at ASC").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func CreateMessagesByList(ctx context.Context, messages []*Message) error {
	err := GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(messages).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func GetMessageIDsByUserIDs(ctx context.Context, userID int64, toUserID int64) ([]*Message, error) {
	res := make([]*Message, 0)
	if err := GetDB().WithContext(ctx).Select("id").Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)", userID, toUserID, toUserID, userID).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func GetMessageByID(ctx context.Context, messageID int64) (*Message, error) {
	res := new(Message)
	if err := GetDB().WithContext(ctx).Select("id, from_user_id, to_user_id, content, created_at").Where("id = ?", messageID).First(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func GetFriendLatestMessage(ctx context.Context, userID int64, toUserID int64) (*Message, error) {
	var res *Message
	if err := GetDB().WithContext(ctx).Select("id, from_user_id, to_user_id, content, created_at").Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)", userID, toUserID, toUserID, userID).Order("created_at DESC").Limit(1).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
