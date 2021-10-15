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
// @Param Body body rpc.CreateChargePointReq true "Post CreateChargePointReq body"
// @Success 200 {object} rpc.CreateChargePointResp
// @Router /v2/ocpp/charge_points [post]
func (api *ChargePointsAPI) CreateChargePoint(c *gin.Context) {
	req := &rpc.CreateChargePointReq{}
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

	res, err := api.chargepointService.CreateChargePoint(ctx, req)

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
// @Param Body body rpc.CreateChargePointIdTagReq true "Post CreateChargePointIdTagReq body"
// @Success 200 {object} rpc.CreateChargePointIdTagReq
// @Router /v2/ocpp/charge_points/id_tags [post]
func (api *ChargePointsAPI) CreateChargePointIdTag(c *gin.Context) {
	req := &rpc.CreateChargePointIdTagReq{}
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

	res, err := api.chargepointService.CreateChargePointIdTag(ctx, req)

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
