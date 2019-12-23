package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/metegol-project/services"
	"github.com/metegol-project/views"
	"net/http"
)

type Controller struct {
	Service *services.Service
}

func (c *Controller) AddUsers(ctx *gin.Context) {
	var request views.AddUsersRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := c.Service.AddUsers(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *Controller) GetUsers(ctx *gin.Context) {
	result, err := c.Service.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *Controller) GetMatches(ctx *gin.Context) {
	tournament := ctx.Param("tournament")

	if tournament == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "tournament cannot be empty",
		})
		return
	}

	result, err := c.Service.GetMatches(tournament)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *Controller) PlayMatch(ctx *gin.Context) {
	var request views.PlayMatchRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := c.Service.PlayMatch(request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *Controller) GetScores(ctx *gin.Context) {
	tournament := ctx.Param("tournament")

	if tournament == "" {
		panic("tournament cannot be empty")
	}

	ctx.JSON(http.StatusOK, c.Service.GetScores(tournament))
}

func (c *Controller) WipeData(ctx *gin.Context) {
	tournament := ctx.Param("tournament")

	if tournament == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "tournament cannot be empty",
		})
		return
	}

	ctx.JSON(http.StatusOK, c.Service.WipeData(tournament))
}
