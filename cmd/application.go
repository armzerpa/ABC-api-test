package cmd

import (
	"database/sql"
	"log"

	"github.com/armzerpa/ABC-api-test/cmd/api/handler"
	"github.com/armzerpa/ABC-api-test/cmd/api/repository/db"
	"github.com/armzerpa/ABC-api-test/cmd/config"
	"github.com/gin-gonic/gin"
)

type App struct {
	Router *gin.Engine
}

func (a *App) Initialize(config *config.Config) {
	db := a.initDB(config.DB)
	bookHandler := handler.NewBookHandler(db)
	a.initRouter(config.Route, bookHandler)
}

func (a *App) initDB(config *config.DBConfig) *sql.DB {
	db, err := db.InitDatabase(*config)
	if err != nil {
		log.Println("Could not connect to database")
	}
	return db
}

func (a *App) initRouter(config *config.RouteConfig, handlerBook *handler.HandlerBook) {
	a.Router = gin.Default()
	version := a.Router.Group(config.Version)
	{
		version.GET("/ping", handler.Ping)
		version.GET("/user", handlerBook.GetUsers)
		version.GET("/user/:id", handlerBook.GetUserById)
		version.DELETE("/user/:id", handlerBook.DeleteUserById)
		version.POST("/user", handlerBook.CreateUser)
	}
	a.Router.Run(config.Port)
}
