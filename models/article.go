package models

import (
	"context"
	"time"
)

type Article struct {
	ID        int       `gorm:"AUTO_INCREMENT" json:"id"`
	Author    string    `gorm:"type:text" json:"author" form:"author" validate:"required"`
	Title     string    `gorm:"type:text" json:"title" form:"title" validate:"required"`
	Body      string    `gorm:"type:text" json:"body" form:"body" validate:"required"`
	CreatedAt time.Time `gorm:"timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
}

type ArticleService interface {
	Store(context.Context, *Article) error
	GetArticle(ctx context.Context, author string, query string) ([]Article, error)
}

type ArticleRepository interface {
	Store(ctx context.Context, u *Article) error
	GetArticle(ctx context.Context, author string, query string) ([]Article, error)
}
