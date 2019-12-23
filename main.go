package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/metegol-project/controllers"
	"github.com/metegol-project/models"
	"github.com/metegol-project/services"

	"net/http"
)

func main() {
	/*
		[] Response Views
		[] Logging
		[] Error Handling
		[] Environment Vars Handling (dev/prod)
		[] Routes
		[] Special Business Rules
		[] Schema Improvement (constraints)
		[] Authentication
		[] Atomic Operations
		[] Health Check
		[] Process Teardown
		[] Migrations Handling
		[] Unit Testing
		[] Integration Testing
		[] Docker
		[] CI
		[] CD
	*/

	db, err := gorm.Open("mysql", "root:root@(localhost:3306)/metegol_db?parseTime=true")
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.Match{}, &models.User{})

	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})

	controller := controllers.Controller{
		Service: services.NewService(db),
	}

	group := r.Group("/metegol")
	group.POST("/users", controller.AddUsers)
	group.GET("/users", controller.GetUsers)
	group.GET("/matches/:tournament", controller.GetMatches)
	group.PUT("/matches", controller.PlayMatch)
	group.GET("/scores/:tournament", controller.GetScores)
	group.DELETE("/data/:tournament", controller.WipeData)

	if err := r.Run(); err != nil {
		panic(err)
	}
}
