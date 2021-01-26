package main

import (
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout)
	logger.Info().Msg(time.Now().Local().String() + "-> App starts, env = " + os.Getenv("ENV"))
	log.Println()
	// app := server.App{}
	// app.InitializeApp()
	// app.Run(":" + os.Getenv("PORT"))
}
