package main

import (
	"log"
	"time"

	"github.com/ashish9868/meracloud/controllers"
	"github.com/ashish9868/meracloud/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	// swagger embed files
)

func main() {

	location, _ := time.LoadLocation("Asia/Kolkata")
	time.Local = location

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		userController := controllers.UserController{}
		v1.POST("/register", userController.Register)
		v1.POST("/login", userController.Login)
		v1.GET("/logout", middleware.BearerAuth(), userController.Logout)
		v1.GET("/download/:id", middleware.BearerAuth(), controllers.UploadController{}.DownloadFile)

		// uploads controller
		uploadsGroup := v1.Group("/uploads").Use(middleware.BearerAuth())
		{
			uploadController := controllers.UploadController{}
			uploadsGroup.GET("/", uploadController.GetFileInfo)
			uploadsGroup.GET("/:id", uploadController.GetFileInfo)
			uploadsGroup.POST("/", uploadController.Upload)
		}
	}

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
