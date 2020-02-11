package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func InitDB() *gorm.DB {
	var err error
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
		os.Getenv("MYSQL_USERNAME"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("DATABASE_NAME"),
	)
	db, errOpen := gorm.Open("mysql", dsn)
	if errOpen != nil {
		log.Fatal(errOpen.Error())
	}
	db.SingularTable(true)

	err = db.DB().Ping()

	if err != nil {
		log.Panic(err)
		fmt.Print(err.Error())
	} else {
		fmt.Printf("Successfully connected to MySQL")
	}
	db.DB().SetMaxIdleConns(60)
	//db.LogMode(true)
	//db.DB().SetMaxOpenConns(80)

	return db
}
