package env

import (
	"github.com/joho/godotenv"
)

func MustLoad() {
	err := godotenv.Load()
	if err != nil {
		panic(".env file does not exists")
	}
}
