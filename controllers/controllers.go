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
		panic(err)
	}

	result := c.Service.AddUsers(request)

	ctx.JSON(http.StatusOK, result)
}

func (c *Controller) GetUsers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.Service.GetUsers())
}

func (c *Controller) GetMatches(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.Service.GetMatches())
}

func (c *Controller) PlayMatch(ctx *gin.Context) {
	var request views.PlayMatchRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, c.Service.PlayMatch(request))
}

func (c *Controller) GetScores(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.Service.GetScores())
}
