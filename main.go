package main

import (
	"fmt"
	"log"
	"time"

	"github.com/leopurba/go-article/config"
	"github.com/leopurba/go-article/controller"
	"github.com/leopurba/go-article/database"
	"github.com/leopurba/go-article/repository"
	"github.com/leopurba/go-article/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	dbConn, err := database.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	//migration
	nM := database.NewMigration(dbConn)
	err = nM.Migration()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		db, err := dbConn.DB()
		if err != nil {
			log.Fatal(err)
		}
		err = db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("redis")
	redisClient, err := database.NewRClient()
	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()
	fmt.Println("redis done")
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	timeoutContext := time.Duration(config.Cfg().Timeout) * time.Second

	articleRepository := repository.NewArticleRepository(dbConn, redisClient)

	articleService := service.NewArticleService(articleRepository, timeoutContext)
	controller.NewArticleHandler(e, articleService)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.Cfg().AppPort)))
}
