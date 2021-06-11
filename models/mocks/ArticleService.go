// Code generated by mockery v1.0.0. DO NOT EDIT.
package mocks

import (
	"context"

	"github.com/leopurba/go-article/models"
	"github.com/stretchr/testify/mock"
)

type ArticleService struct {
	mock.Mock
}

func (_m *ArticleService) Store(_a0 context.Context, _a1 *models.Article) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Article) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *ArticleService) GetArticle(ctx context.Context, author string, query string) ([]models.Article, error) {
	ret := _m.Called(ctx, author, query)

	var r0 []models.Article
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []models.Article); ok {
		r0 = rf(ctx, author, query)
	} else {
		r0 = ret.Get(0).([]models.Article)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, author, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}