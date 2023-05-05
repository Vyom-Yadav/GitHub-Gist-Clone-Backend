package main

import (
	"fmt"
	"log"

	"github.com/Vyom-Yadav/GitHub-Gist-Clone/initializers"
	"github.com/Vyom-Yadav/GitHub-Gist-Clone/models"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
}

func main() {
	err := initializers.DB.AutoMigrate(&models.User{}, &models.UserMetadata{}, &models.Gist{}, &models.Comment{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("? Migration complete")
}
