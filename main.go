package main

import (
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/unnamedxaer/gymm-api/server"
)

func main() {
	log.Println(time.Now().Local().String() + "-> App starts, env = ")
	app := server.App{}
	app.InitializeApp()
	app.Run(":" + os.Getenv("PORT"))

}
