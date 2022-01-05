package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Beep-Technologies/beepbeep3-iam/pkg/constants"
	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	chargepoint "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/charge_point"
	statusnotification "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/status_notification"
)

type ChargePointsAPI struct {
	chargepointService        *chargepoint.Service
	statusnotificationService *statusnotification.Service
}

func NewChargePointsAPI(cpS *chargepoint.Service, snS *statusnotification.Service) *ChargePointsAPI {
	return &ChargePointsAPI{
		chargepointService:        cpS,
		statusnotificationService: snS,
	}
}

// GetChargePoint gets the status of the charge point
// @Summary Get charge point
// @Tags Charge Points
// @Produce json
// @Param entityCode   path string true "entity code"
// @Param chargePointIdentifier path string true "charge point identifier"
// @Success 200 {object} rpc.GetChargePointResp
// @Router /v2/ocpp/charge-points/{entityCode}/{chargePointIdentifier} [get]
func (api *ChargePointsAPI) GetChargePoint(c *gin.Context) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.CtxKey("gin"), c)

	// get the entity code
	entityCode := c.Param("entityCode")

	// get the charge point identifier and decode it
	// charge point identifiers are percent-encoded
	chargePointIdentifier := c.Param("chargePointIdentifier")

	req := &rpc.GetChargePointReq{
		EntityCode:            entityCode,
		ChargePointIdentifier: chargePointIdentifier,
	}

	res, err := api.chargepointService.GetChargePoint(ctx, req)

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
		"data":    res.ChargePoint,
	})
}

// GetChargePointConnectorStatus gets the latest status of the connectors on the charge point
// @Summary Get charge point connector status
// @Tags Charge Points
// @Produce json
// @Param entityCode   path string true "entity code"
// @Param chargePointIdentifier path string true "charge point identifier"
// @Success 200 {object} rpc.GetLatestStatusNotificationsResp
// @Router /v2/ocpp/charge-points/{entityCode}/{chargePointIdentifier}/status [get]
func (api *ChargePointsAPI) GetChargePointConnectorStatus(c *gin.Context) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.CtxKey("gin"), c)

	// get the entity code
	entityCode := c.Param("entityCode")

	// get the charge point identifier and decode it
	// charge point identifiers are percent-encoded
	chargePointIdentifier := c.Param("chargePointIdentifier")

	req := &rpc.GetLatestStatusNotificationsReq{
		EntityCode:            entityCode,
		ChargePointIdentifier: chargePointIdentifier,
	}

	res, err := api.statusnotificationService.GetLatestStatusNotifications(ctx, req)

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
