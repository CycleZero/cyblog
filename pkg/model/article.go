package model

import "gorm.io/gorm"

// 文章状态常量
const (
	ArticleStatusDraft   = 1 // 草稿
	ArticleStatusPublish = 2 // 已发布
	ArticleStatusArchive = 3 // 已归档
)

// Article 文章模型
type Article struct {
	gorm.Model
	Title       string `gorm:"size:255;not null;comment:文章标题"`
	Slug        string `gorm:"size:255;unique;comment:文章别名（用于URL）"`
	Content     string `gorm:"type:text;not null;comment:文章内容（Markdown）"`
	HtmlContent string `gorm:"type:text;not null;comment:渲染后的HTML内容"`
	Summary     string `gorm:"size:500;comment:文章摘要"`
	CoverImage  string `gorm:"size:255;comment:封面图URL"`
	Status      int    `gorm:"default:1;comment:状态：1草稿 2已发布 3已归档"`
	Views       int    `gorm:"default:0;comment:浏览量"`
	Likes       int    `gorm:"default:0;comment:点赞数"`
	IsTop       bool   `gorm:"default:false;comment:是否置顶"`
	IsOriginal  bool   `gorm:"default:true;comment:是否原创"`
	AuthorID    uint   `gorm:"not null;comment:作者ID"`
	CategoryID  uint   `gorm:"comment:分类ID"`

	// 关联
	Author   *User     `gorm:"foreignKey:AuthorID"`
	Category *Category `gorm:"foreignKey:CategoryID"`
	Tags     []*Tag    `gorm:"many2many:article_tags;"`
}

// Category 分类模型
type Category struct {
	gorm.Model
	Name        string `gorm:"size:50;not null;unique;comment:分类名称"`
	Slug        string `gorm:"size:50;unique;comment:分类别名"`
	Description string `gorm:"size:200;comment:分类描述"`
	ParentID    uint   `gorm:"default:0;comment:父分类ID"`
	Sort        int    `gorm:"default:0;comment:排序"`

	// 关联
	Articles []*Article `gorm:"foreignKey:CategoryID"`
}

// Tag 标签模型
type Tag struct {
	gorm.Model
	Name  string `gorm:"size:50;not null;unique;comment:标签名称"`
	Slug  string `gorm:"size:50;unique;comment:标签别名"`
	Color string `gorm:"size:7;default:'#165DFF';comment:标签颜色"`
	Count int    `gorm:"default:0;comment:关联文章数"`

	// 关联
	Articles []*Article `gorm:"many2many:article_tags;"`
}

// ArticleTag 文章标签关联表模型
type ArticleTag struct {
	ArticleID uint `gorm:"primaryKey"`
	TagID     uint `gorm:"primaryKey"`
	CreatedAt gorm.DeletedAt
}
