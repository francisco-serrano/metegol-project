package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/metegol-project/services"
	"github.com/metegol-project/utils"
	"github.com/metegol-project/viewmodels"
	"net/http"
)

type Controller struct {
	Service *services.Service
}

func (c *Controller) AddUsers(ctx *gin.Context) {
	var request viewmodels.AddUsersRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.SetErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	result, err := c.Service.AddUsers(request)
	if err != nil {
		utils.SetErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	utils.SetResponse(ctx, http.StatusCreated, result)
}

func (c *Controller) GetUsers(ctx *gin.Context) {
	result, err := c.Service.GetUsers()
	if err != nil {
		utils.SetErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	utils.SetResponse(ctx, http.StatusOK, result)
}

func (c *Controller) GetMatches(ctx *gin.Context) {
	tournament := ctx.Param("tournament")

	if tournament == "" {
		err := errors.New("tournament cannot be empty")
		utils.SetErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	result, err := c.Service.GetMatches(tournament)
	if err != nil {
		utils.SetErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	utils.SetResponse(ctx, http.StatusOK, result)
}

func (c *Controller) PlayMatch(ctx *gin.Context) {
	var request viewmodels.PlayMatchRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.SetErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	result, err := c.Service.PlayMatch(request)
	if err != nil {
		utils.SetErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	utils.SetResponse(ctx, http.StatusOK, result)
}

func (c *Controller) GetScores(ctx *gin.Context) {
	tournament := ctx.Param("tournament")
	if tournament == "" {
		err := errors.New("tournament cannot be empty")
		utils.SetErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	result, err := c.Service.GetScores(tournament)
	if err != nil {
		utils.SetErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	utils.SetResponse(ctx, http.StatusOK, result)
}

func (c *Controller) WipeData(ctx *gin.Context) {
	tournament := ctx.Param("tournament")
	if tournament == "" {
		err := errors.New("tournament cannot be empty")
		utils.SetErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	result := c.Service.WipeData(tournament)

	utils.SetResponse(ctx, http.StatusOK, result)
}
