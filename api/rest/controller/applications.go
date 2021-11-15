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

// SetApplicationCallback sets the callback url for an application
// @Summary Set Callback URL for Application
// @Tags Applications
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Body body rpc.CreateApplicationCallbackReqPublic true "Post CreateApplicationCallbackReq body"
// @Success 200 {object} rpc.CreateApplicationCallbackResp
// @Router /v2/ocpp/applications/callbacks [post]
func (api *ApplicationsAPI) SetApplicationCallback(c *gin.Context) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.CtxKey("gin"), c)

	applicationId := c.GetString("application_id")

	req := &rpc.CreateApplicationCallbackReqPublic{}
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

	res, err := api.applicationService.SetApplicationCallback(ctx, &rpc.CreateApplicationCallbackReq{
		ApplicationId: applicationId,
		CallbackEvent: req.CallbackEvent,
		CallbackUrl:   req.CallbackUrl,
	})

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

// GetApplicationCallbacks gets the callback urls for an application
// @Summary Get Callback URL for Application
// @Tags Applications
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} rpc.GetApplicationCallbacksResp
// @Router /v2/ocpp/applications/callbacks [get]
func (api *ApplicationsAPI) GetApplicationCallbacks(c *gin.Context) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.CtxKey("gin"), c)

	applicationId := c.GetString("application_id")

	res, err := api.applicationService.GetApplicationCallbacks(ctx, &rpc.GetApplicationCallbacksReq{
		ApplicationId: applicationId,
	})
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
