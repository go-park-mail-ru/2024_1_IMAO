package utils

import (
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	const projectDirName = "2024_1_IMAO"
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	err := godotenv.Load(string(rootPath) + `/.env`)

	if err != nil {
		log.Println("Error loading .env file")
	}
}
