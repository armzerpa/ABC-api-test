package cmd

import (
	"database/sql"

	"github.com/armzerpa/ABC-api-test/cmd/api/handler"
	"github.com/armzerpa/ABC-api-test/cmd/api/repository/db"
	"github.com/armzerpa/ABC-api-test/cmd/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type App struct {
	Router *gin.Engine
}

var log = logrus.New()

func (a *App) Initialize(config *config.Config) {
	log.WithFields(logrus.Fields{}).Info("starting api service")
	db := a.initDB(config.DB)
	userHandler := handler.NewUserHandler(db)
	fileHandler := handler.NewFileHandler(db)
	a.initRouter(config.Route, userHandler, fileHandler)
}

func (a *App) initDB(config *config.DBConfig) *sql.DB {
	db, err := db.InitDatabase(*config)
	if err != nil {
		log.WithFields(logrus.Fields{}).Error("couldnot connect with database")
	}
	return db
}

func (a *App) initRouter(config *config.RouteConfig, handlerUser *handler.HandlerUser, handlerFile *handler.HandlerFile) {
	a.Router = gin.Default()
	authorized := a.Router.Group(config.Version, gin.BasicAuth(gin.Accounts{
		"user": "1234",
	}))
	authorized.GET("/ping", handler.Ping)
	authorized.GET("/user", handlerUser.GetUsers)
	authorized.GET("/user/:id", handlerUser.GetUserById)
	authorized.DELETE("/user/:id", handlerUser.DeleteUserById)
	authorized.POST("/user", handlerUser.CreateUser)
	authorized.PUT("/user/:id", handlerUser.UpdateUser)

	authorized.GET("/user/:id/file", handlerFile.GetByUserId)
	authorized.DELETE("/user/:id/file", handlerFile.DeleteByUserId)
	authorized.POST("/user/:id/file", handlerFile.CreateFile)

	a.Router.Run(config.Port)
}
