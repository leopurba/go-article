package controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/bxcodec/faker"
	"github.com/labstack/echo/v4"
	"github.com/leopurba/go-article/controller"
	"github.com/leopurba/go-article/models"
	"github.com/leopurba/go-article/models/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	mockArticle := models.Article{
		Author:    "Author",
		Title:     "Title",
		Body:      "Body",
		CreatedAt: time.Now(),
	}

	tempMockArticle := mockArticle
	tempMockArticle.ID = 0
	mockService := new(mocks.ArticleService)

	j, err := json.Marshal(tempMockArticle)
	assert.NoError(t, err)

	mockService.On("Store", mock.Anything, mock.AnythingOfType("*models.Article")).Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/articles", strings.NewReader(string(j)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/article")

	handler := controller.ArticleHandler{
		ArticleService: mockService,
	}
	err = handler.Store(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockService.AssertExpectations(t)
}

func TestGetArticle(t *testing.T) {
	var mockArticle models.Article
	err := faker.FakeData(&mockArticle)
	assert.NoError(t, err)
	mockService := new(mocks.ArticleService)
	mockListArticle := make([]models.Article, 0)
	mockListArticle = append(mockListArticle, mockArticle)
	author := ""
	query := ""
	mockService.On("GetArticle", mock.Anything, author, query).Return(mockListArticle, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/article?author="+author+"&query="+query, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := controller.ArticleHandler{
		ArticleService: mockService,
	}
	err = handler.GetArticle(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestGetArticleError(t *testing.T) {
	mockService := new(mocks.ArticleService)
	author := ""
	query := ""
	mockService.On("GetArticle", mock.Anything, author, query).Return([]models.Article{}, models.ErrInternalServerError)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/article?author="+author+"&query="+query, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := controller.ArticleHandler{
		ArticleService: mockService,
	}
	err = handler.GetArticle(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockService.AssertExpectations(t)
}
