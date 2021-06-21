package repository

import (
	"context"
	"fmt"
	"strings"

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
	if create := ar.Conn.WithContext(ctx).Create(&article); create.RowsAffected == 0 {
		return models.ErrInternalServerError
	}
	iter := ar.RClient.Conn().Scan(ctx, 0, "article_*", 0).Iterator()
	for iter.Next(ctx) {
		key := strings.Split(iter.Val(), "_")
		if key[1] != "" {
			if strings.Contains(article.Author, key[1]) {
				if key[2] != "" {
					if strings.Contains(article.Title, key[2]) || strings.Contains(article.Body, key[2]) {
						var cArticel []models.Article
						err = ar.RClient.Cache().Get(ctx, iter.Val(), &cArticel)
						if err != nil && err != cache.ErrCacheMiss {
							return err
						}
						err = ar.RClient.Cache().Delete(ctx, iter.Val())
						if err != nil && err != cache.ErrCacheMiss {
							return err
						}
						var preArticle []models.Article
						preArticle = append(preArticle, *article)
						cArticel = append(preArticle, cArticel...)
						ar.RClient.Cache().Set(&cache.Item{
							Ctx:   ctx,
							Key:   iter.Val(),
							Value: cArticel,
							TTL:   config.Cfg().RedisTTL,
						})
					}
				} else {
					var cArticel []models.Article
					err = ar.RClient.Cache().Get(ctx, iter.Val(), &cArticel)
					if err != nil && err != cache.ErrCacheMiss {
						return err
					}
					err = ar.RClient.Cache().Delete(ctx, iter.Val())
					if err != nil && err != cache.ErrCacheMiss {
						return err
					}
					var preArticle []models.Article
					preArticle = append(preArticle, *article)
					cArticel = append(preArticle, cArticel...)
					ar.RClient.Cache().Set(&cache.Item{
						Ctx:   ctx,
						Key:   iter.Val(),
						Value: cArticel,
						TTL:   config.Cfg().RedisTTL,
					})
				}
			}
		} else {
			if key[2] != "" {
				if strings.Contains(article.Title, key[2]) || strings.Contains(article.Body, key[2]) {
					var cArticel []models.Article
					err = ar.RClient.Cache().Get(ctx, iter.Val(), &cArticel)
					if err != nil && err != cache.ErrCacheMiss {
						return err
					}
					err = ar.RClient.Cache().Delete(ctx, iter.Val())
					if err != nil && err != cache.ErrCacheMiss {
						return err
					}
					var preArticle []models.Article
					preArticle = append(preArticle, *article)
					cArticel = append(preArticle, cArticel...)
					ar.RClient.Cache().Set(&cache.Item{
						Ctx:   ctx,
						Key:   iter.Val(),
						Value: cArticel,
						TTL:   config.Cfg().RedisTTL,
					})
				}
			} else {
				var cArticel []models.Article
				err = ar.RClient.Cache().Get(ctx, iter.Val(), &cArticel)
				if err != nil && err != cache.ErrCacheMiss {
					return err
				}
				err = ar.RClient.Cache().Delete(ctx, iter.Val())
				if err != nil && err != cache.ErrCacheMiss {
					return err
				}
				var preArticle []models.Article
				preArticle = append(preArticle, *article)
				cArticel = append(preArticle, cArticel...)
				ar.RClient.Cache().Set(&cache.Item{
					Ctx:   ctx,
					Key:   iter.Val(),
					Value: cArticel,
					TTL:   config.Cfg().RedisTTL,
				})
			}
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
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
