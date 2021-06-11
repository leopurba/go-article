package service

import (
	"context"
	"time"

	"github.com/leopurba/go-article/models"
)

type articleService struct {
	articleRepository models.ArticleRepository
	contextTimeout    time.Duration
}

func NewArticleService(a models.ArticleRepository, timeout time.Duration) models.ArticleService {
	return &articleService{
		articleRepository: a,
		contextTimeout:    timeout,
	}
}

func (a *articleService) Store(c context.Context, u *models.Article) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	err = a.articleRepository.Store(ctx, u)
	if err != nil {
		return
	}
	return
}

func (a *articleService) GetArticle(c context.Context, author string, query string) ([]models.Article, error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err := a.articleRepository.GetArticle(ctx, author, query)
	if err != nil {
		return nil, err
	}

	return res, nil
}
