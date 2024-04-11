package env

import (
	"log"

	"github.com/joho/godotenv"
)

func MustLoad() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file does not exists")
	}
}
