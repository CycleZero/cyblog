package model

import "gorm.io/gorm"

// Comment 评论模型
type Comment struct {
	gorm.Model
	ArticleID uint   `gorm:"not null;index;comment:文章ID"`
	UserID    uint   `gorm:"not null;index;comment:评论者ID"`
	ParentID  uint   `gorm:"default:0;index;comment:父评论ID（0表示一级评论）"`
	Content   string `gorm:"type:text;not null;comment:评论内容"`
	Likes     int    `gorm:"default:0;comment:点赞数"`
	IP        string `gorm:"size:45;comment:评论者IP"`
	UserAgent string `gorm:"size:500;comment:评论者UserAgent"`

	// 关联
	Article *Article   `gorm:"foreignKey:ArticleID"`
	User    *User      `gorm:"foreignKey:UserID"`
	Parent  *Comment   `gorm:"foreignKey:ParentID"`
	Replies []*Comment `gorm:"foreignKey:ParentID;order:created_at asc"`
}

// CommentLike 评论点赞记录模型
type CommentLike struct {
	gorm.Model
	CommentID uint   `gorm:"not null;index;comment:评论ID"`
	UserID    uint   `gorm:"not null;index;comment:用户ID"`
	IP        string `gorm:"size:45;comment:点赞者IP"`
	UserAgent string `gorm:"size:500;comment:点赞者UserAgent"`

	// 关联
	Comment *Comment `gorm:"foreignKey:CommentID"`
	User    *User    `gorm:"foreignKey:UserID"`
}
