package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	service Service
}

func (c *Controller) AddUsers(ctx *gin.Context) {
	var request AddUsersRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		panic(err)
	}

	result := c.service.AddUsers(request)

	ctx.JSON(http.StatusOK, result)
}

func (c *Controller) GetUsers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.service.GetUsers())
}

func (c *Controller) GetMatches(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.service.GetMatches())
}
