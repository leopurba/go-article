package repository

import (
	"context"
	"fmt"

	"github.com/go-redis/cache/v8"
	"github.com/leopurba/go-article/config"
	"github.com/leopurba/go-article/database"
	"github.com/leopurba/go-article/models"

	"gorm.io/gorm"
)

type articleRepository struct {
	Conn    *gorm.DB
	RClient database.RClient
}

func NewArticleRepository(Conn *gorm.DB, RClient database.RClient) models.ArticleRepository {
	return &articleRepository{Conn, RClient}
}

func (ar *articleRepository) Store(ctx context.Context, article *models.Article) (err error) {
	if create := ar.Conn.WithContext(ctx).Create(&article); create.RowsAffected > 0 {
		return nil
	}
	return models.ErrInternalServerError
}

func (ar *articleRepository) GetArticle(ctx context.Context, author string, query string) ([]models.Article, error) {
	var article []models.Article
	err := ar.RClient.Cache().Get(ctx, fmt.Sprintf("article_%s_%s", author, query), &article)
	if err != nil && err != cache.ErrCacheMiss {
		return nil, err
	} else if err == nil {
		return article, nil
	}

	if err := ar.Conn.WithContext(ctx).Where(ar.Conn.WithContext(ctx).Where("title LIKE ?", "%"+query+"%").Or("body LIKE ?", "%"+query+"%")).Where("author LIKE ?", "%"+author+"%").Order("id desc").Find(&article).Error; err != nil {
		return nil, err
	}
	if len(article) == 0 {
		return article, nil
	}
	return article, ar.RClient.Cache().Set(&cache.Item{
		Ctx:   ctx,
		Key:   fmt.Sprintf("article_%s_%s", author, query),
		Value: article,
		TTL:   config.Cfg().RedisTTL,
	})
}
