package main

import (
	"net/http"
	"tusk/config"
	"tusk/controllers"
	"tusk/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// Databse
	db := config.DatabaseConnection()
	db.AutoMigrate(&models.User{}, &models.Task{})
	config.CreateOwnerAccount(db)

	//controller
	UserController := controllers.UserController{DB: db}
	TaskController := controllers.TaskController{DB: db}

	// router
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Welcome Golang")
	})

	router.POST("/users/login", UserController.Login)
	router.POST("/users", UserController.CreateAccount)
	router.DELETE("/users/:id", UserController.Delete)
	router.GET("/users/employee", UserController.GetEmployee)

	router.POST("/tasks", TaskController.Create)
	router.DELETE("/tasks/:id", TaskController.Delete)
	router.PATCH("/tasks/:id/submit", TaskController.Submit)
	router.PATCH("/tasks/:id/reject", TaskController.Reject)
	router.PATCH("/tasks/:id/fix", TaskController.Fix)
	router.PATCH("/tasks/:id/approve", TaskController.Approve)
	router.GET("/tasks/:id", TaskController.FindById)
	router.GET("/tasks/review/asc", TaskController.NeedToBeReview)
	router.GET("/tasks/progress/:userId", TaskController.ProgressTasks)
	router.GET("/tasks/stat/:userId", TaskController.Statistic)
	router.GET("/tasks/user/:userId/:status", TaskController.FindByUserAndStatus)

	router.Static("/attachments", "./attachments")

	router.Run("0.0.0.0:8080")
}
