package main

import (
	"github.com/gin-gonic/gin"
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
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})

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

	Run()
}
