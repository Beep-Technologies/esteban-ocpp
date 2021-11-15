package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Beep-Technologies/beepbeep3-iam/pkg/constants"
	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	chargepoint "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/charge_point"
)

type ChargePointsAPI struct {
	chargepointService *chargepoint.Service
}

func NewChargePointsAPI(cS *chargepoint.Service) *ChargePointsAPI {
	return &ChargePointsAPI{
		chargepointService: cS,
	}
}

// CreateChargePoint creates a charge point
// @Summary Create a Charge Point
// @Tags Charge Points
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Body body rpc.CreateChargePointReqPublic true "Post CreateChargePointReq body"
// @Success 200 {object} rpc.CreateChargePointResp
// @Router /v2/ocpp/charge_points [post]
func (api *ChargePointsAPI) CreateChargePoint(c *gin.Context) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.CtxKey("gin"), c)

	applicationId := c.GetString("application_id")

	req := &rpc.CreateChargePointReqPublic{}
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

	res, err := api.chargepointService.CreateChargePoint(ctx, &rpc.CreateChargePointReq{
		ApplicationId:         applicationId,
		ChargePointIdentifier: req.ChargePointIdentifier,
		EntityCode:            req.EntityCode,
		OcppProtocol:          req.OcppProtocol,
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

// CreateChargePointIdTag creates an id tag associated with a charge point
// @Summary Create a Charge Point ID tag
// @Tags Charge Points
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Body body rpc.CreateChargePointIdTagReqPublic true "Post CreateChargePointIdTagReq body"
// @Success 200 {object} rpc.CreateChargePointIdTagReq
// @Router /v2/ocpp/charge_points/id_tags [post]
func (api *ChargePointsAPI) CreateChargePointIdTag(c *gin.Context) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.CtxKey("gin"), c)

	applicationId := c.GetString("application_id")

	req := &rpc.CreateChargePointIdTagReqPublic{}
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

	res, err := api.chargepointService.CreateChargePointIdTag(ctx, &rpc.CreateChargePointIdTagReq{
		ApplicationId:         applicationId,
		ChargePointIdentifier: req.ChargePointIdentifier,
		IdTag:                 req.IdTag,
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
