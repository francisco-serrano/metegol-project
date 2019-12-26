package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/metegol-project/controllers"
	"github.com/metegol-project/models"
	"github.com/metegol-project/services"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func checkEnvironmentVariables() error {
	envVars := []string{"DB_USER", "DB_PASS", "DB_HOST", "DB_PORT", "DB_NAME"}

	for _, v := range envVars {
		if myVar := os.Getenv(v); myVar == "" {
			return errors.New(fmt.Sprintf("%s not provided", v))
		}
	}

	return nil
}

func InitializeDatabase() (*gorm.DB, error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dbConnection := fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true", user, password, host, port, dbName)

	db, err := gorm.Open("mysql", dbConnection)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Match{}, &models.User{})

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(time.Hour)

	return db, nil
}

func InitializeLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Out = os.Stdout
	logger.Level = logrus.DebugLevel
	logger.Formatter = &logrus.JSONFormatter{}

	return logger
}

func InitializeRoutes(engine *gin.Engine, db *gorm.DB, logger *logrus.Logger) error {
	service, err := services.NewService(db, logger)
	if err != nil {
		return err
	}

	controller := controllers.Controller{
		Service: service,
	}

	group := engine.Group("/metegol")
	group.POST("/users", controller.AddUsers)
	group.GET("/users", controller.GetUsers)
	group.GET("/matches/:tournament", controller.GetMatches)
	group.PUT("/matches", controller.PlayMatch)
	group.GET("/scores/:tournament", controller.GetScores)
	group.DELETE("/data/:tournament", controller.WipeData)

	return nil
}

func main() {
	if err := checkEnvironmentVariables(); err != nil {
		panic(err)
	}

	db, err := InitializeDatabase()
	if err != nil {
		panic(err)
	}

	logger := InitializeLogger()

	r := gin.Default()

	if err := InitializeRoutes(r, db, logger); err != nil {
		panic(err)
	}

	if err := r.Run(); err != nil {
		panic(err)
	}
}
