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
	userHandler := handler.NewUserHandler(db)
	fileHandler := handler.NewFileHandler(db)
	a.initRouter(config.Route, userHandler, fileHandler)
}

func (a *App) initDB(config *config.DBConfig) *sql.DB {
	db, err := db.InitDatabase(*config)
	if err != nil {
		log.Println("Could not connect to database")
	}
	return db
}

func (a *App) initRouter(config *config.RouteConfig, handlerUser *handler.HandlerUser, handlerFile *handler.HandlerFile) {
	a.Router = gin.Default()
	version := a.Router.Group(config.Version)
	{
		version.GET("/ping", handler.Ping)
		version.GET("/user", handlerUser.GetUsers)
		version.GET("/user/:id", handlerUser.GetUserById)
		version.DELETE("/user/:id", handlerUser.DeleteUserById)
		version.POST("/user", handlerUser.CreateUser)
		version.PUT("/user/:id", handlerUser.UpdateUser)

		version.GET("/user/:id/file", handlerFile.GetByUserId)
		version.DELETE("/user/:id/file", handlerFile.DeleteByUserId)
		version.POST("/user/:id/file", handlerFile.CreateFile)
	}
	a.Router.Run(config.Port)
}
