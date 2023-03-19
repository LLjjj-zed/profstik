package mysql

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
