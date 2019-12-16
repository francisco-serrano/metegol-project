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

//func amountMatches(numberPlayers int) int {
//	return int(0.5 * float32(numberPlayers) * (float32(numberPlayers) - 1))
//}
//
//func amountTeams(numberPlayers int) int {
//	return combinatorial(numberPlayers, 2)
//}
//
//func factorial(n int) int {
//	result := n
//	for i := n-1; i > 0; i-- {
//		result *= i
//	}
//
//	return result
//}
//
//func combinatorial(n, k int) int {
//	numerator := factorial(n)
//	denominator := factorial(k) * factorial(n - k)
//
//	return numerator / denominator
//}
//
//func generateTeams(players []string) []string {
//	var teams []string
//	for _, p1 := range players {
//		for _, p2 := range players {
//			if p1 != p2 {
//				teams = append(teams, p1+p2)
//			}
//		}
//	}
//
//	return teams
//}

func main() {
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

	r.POST("/metegol/users", controller.AddUsers)
	r.GET("/metegol/users", controller.GetUsers)
	r.GET("/metegol/matches", controller.GetMatches)
	r.PUT("/metegol/matches", controller.PlayMatch)
	r.GET("/metegol/scores", controller.GetScores)

	if err := r.Run(); err != nil {
		panic(err)
	}

	//numberPlayers := 4
	//numberTeams := amountTeams(numberPlayers)
	//numberMatches := amountMatches(numberTeams)
	//
	//fmt.Println("number of players", numberPlayers)
	//fmt.Println("number of teams", numberTeams)
	//fmt.Println("number of matches", numberMatches)

	//Run()
}
