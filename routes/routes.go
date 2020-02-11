package routes

import (
	"log"
	"os"
	"samplesentry/config"
	slog "samplesentry/sentrylog"
	"samplesentry/users"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(config.MaxAllowed(30))
	gin.SetMode(os.Getenv("GIN_MODE"))
	db := config.InitDB()
	var logger *log.Logger
	sentrylog := slog.New(logger)
	v1 := r.Group("/api")

	users.Users(v1.Group("/users"), db, sentrylog)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
