package users

import (
	"samplesentry/sentrylog"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func Users(router *gin.RouterGroup, db *gorm.DB, logging *sentrylog.Logging) {
	handlers := New(db, logging)
	router.GET("/", handlers.RetrieveUser)
}
