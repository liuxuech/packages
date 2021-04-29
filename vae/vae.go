package main

import (
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	app := gin.New()

	if err := app.Run(":8080"); err != nil {
		os.Exit(1)
	}
}
