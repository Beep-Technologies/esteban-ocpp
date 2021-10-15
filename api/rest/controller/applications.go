package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Beep-Technologies/beepbeep3-iam/pkg/constants"
	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	application "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/application"
)

type ApplicationsAPI struct {
	applicationService *application.Service
}

func NewApplicationsAPI(aS *application.Service) *ApplicationsAPI {
	return &ApplicationsAPI{
		applicationService: aS,
	}
}

// CreateApplication creates an application
// @Summary Create an Application
// @Tags Applications
// @Accept json
// @Produce json
// @Param Body body rpc.CreateApplicationReq true "Post CreateApplicationReq body"
// @Success 200 {object} rpc.CreateApplicationResp
// @Router /v2/ocpp/applications/ [post]
func (api *ApplicationsAPI) CreateApplication(c *gin.Context) {
	req := &rpc.CreateApplicationReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"status":  http.StatusUnprocessableEntity,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.CtxKey("gin"), c)

	res, err := api.applicationService.CreateApplication(ctx, req)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success",
		"data":    res,
	})
}

// SetApplicationCallback sets the callback url for an application
// @Summary Set Callback URL for Application
// @Tags Applications
// @Accept json
// @Produce json
// @Param Body body rpc.CreateApplicationCallbackReq true "Post CreateApplicationCallbackReq body"
// @Success 200 {object} rpc.CreateApplicationCallbackResp
// @Router /v2/ocpp/applications/callbacks [post]
func (api *ApplicationsAPI) SetApplicationCallback(c *gin.Context) {
	req := &rpc.CreateApplicationCallbackReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"status":  http.StatusUnprocessableEntity,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.CtxKey("gin"), c)

	res, err := api.applicationService.SetApplicationCallback(ctx, req)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Success",
		"data":    res,
	})
}
