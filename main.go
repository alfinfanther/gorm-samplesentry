package main

import (
	"log"
	"net/http"
	"os"
	"samplesentry/routes"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	if len(os.Args) > 1 && os.Args[1] == "local" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file not FOUND")
		}
	}

	r := routes.SetupRouter()
	s := &http.Server{
		Addr:           ":8090",
		Handler:        r,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
