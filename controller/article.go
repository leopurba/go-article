package controller

import (
	"net/http"

	"github.com/leopurba/go-article/models"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type ResponseError struct {
	Message string `json:"message"`
}

type ArticleHandler struct {
	ArticleService models.ArticleService
}

func NewArticleHandler(e *echo.Echo, us models.ArticleService) {
	handler := &ArticleHandler{
		ArticleService: us,
	}
	e.POST("/articles", handler.Store)
	e.GET("/articles", handler.GetArticle)
}

func (a *ArticleHandler) Store(c echo.Context) (err error) {
	var article models.Article
	err = c.Bind(&article)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, ResponseError{Message: err.Error()})
	}

	var ok bool
	if ok, err = isRequestValid(&article); !ok {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	ctx := c.Request().Context()
	err = a.ArticleService.Store(ctx, &article)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, article)
}

func (a *ArticleHandler) GetArticle(c echo.Context) (err error) {
	author := c.QueryParam("author")
	query := c.QueryParam("query")

	ctx := c.Request().Context()
	list, err := a.ArticleService.GetArticle(ctx, author, query)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, list)
}

func isRequestValid(m *models.Article) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	switch err {
	case models.ErrInternalServerError:
		return http.StatusInternalServerError
	case models.ErrNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
