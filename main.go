package main

import (
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth/gothic"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	gothic.Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
}

func main() {
	providerIndex := gothSetup()
	db := gormSetup()
	r := ginSetup(db, providerIndex)
	r.Static("/public", "./public")
	r.Run(":8080")
}
