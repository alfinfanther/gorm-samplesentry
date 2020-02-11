package users

import (
	"net/http"
	"samplesentry/models"
	"samplesentry/sentrylog"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Db struct {
	GormDb *gorm.DB
	Log    *sentrylog.Logging
}

func New(db *gorm.DB, logger *sentrylog.Logging) *Db {
	return new(Db).init(db, logger)
}
func (e *Db) init(db *gorm.DB, logger *sentrylog.Logging) *Db {
	e.GormDb = db
	e.Log = logger
	return e
}

func (e *Db) RetrieveUser(c *gin.Context) {
	pagination := models.Pagination{}
	db, err := pagination.Paginate(c, e.GormDb)
	if err != nil {
		e.Log.Error(err.Error)
		c.JSON(400, gin.H{"error": err.Error()})
	}
	var users []models.Users
	args := []interface{}{1}
	conditions := "email != ?"
	errDb := db.Where(conditions, args).Preload("Users").Find(&users)
	// errDb := db.Where(conditions, args).Find(&users)
	mp := map[string]string{
		"query_conditions": conditions,
		"join_conditions":  "",
	}
	pagination.CountRecords("users", mp, args)

	var index int
	if len(users) < 1 {
		index = 0
	} else {
		index = int(users[len(users)-1].Id)
	}
	pagination.SetHeaderLink(c, index)
	var usr []models.Users
	for _, d := range users {
		u := models.Users{
			Id:       d.Id,
			Email:    d.Email,
			Fullname: d.Fullname,
		}
		usr = append(usr, u)
	}

	if errDb.Error != nil {
		e.Log.Error(errDb.Error.Error())
		c.JSON(http.StatusBadRequest, models.APImessage{
			Message: errDb.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, usr)

}
