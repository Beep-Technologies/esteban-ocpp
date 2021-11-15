package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Beep-Technologies/beepbeep3-iam/pkg/constants"
	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/operation"
	statusnotification "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/status_notification"
)

type OperationsAPI struct {
	operationService          *operation.Service
	statusnotificationService *statusnotification.Service
}

func NewOperationsAPI(oS *operation.Service, snS *statusnotification.Service) *OperationsAPI {
	return &OperationsAPI{
		operationService:          oS,
		statusnotificationService: snS,
	}
}

// RemoteStartTransaction has the Central System request the Charge Point to start a transaction
// @Summary Request Charge Point to Start a Transaction
// @Tags Operations
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Body body rpc.RemoteStartTransactionReqPublic true "Post RemoteStartTransactionReq body"
// @Success 200 {object} rpc.RemoteStartTransactionResp
// @Router /v2/ocpp/operations/remote-start-transaction [post]
func (api *OperationsAPI) RemoteStartTransaction(c *gin.Context) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.CtxKey("gin"), c)

	applicationId := c.GetString("application_id")

	req := &rpc.RemoteStartTransactionReqPublic{}
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

	res, err := api.operationService.RemoteStartTransaction(ctx, &rpc.RemoteStartTransactionReq{
		ApplicationId:         applicationId,
		ChargePointIdentifier: req.ChargePointIdentifier,
		ConnectorId:           req.ConnectorId,
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

// RemoteStartTransaction has the Central System request the Charge Point to stop a transaction
// @Summary Request Charge Point to Stop a Transaction
// @Tags Operations
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Body body rpc.RemoteStopTransactionReqPublic true "Post RemoteStopTransactionReq body"
// @Success 200 {object} rpc.RemoteStopTransactionResp
// @Router /v2/ocpp/operations/remote-stop-transaction [post]
func (api *OperationsAPI) RemoteStopTransaction(c *gin.Context) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.CtxKey("gin"), c)

	applicationId := c.GetString("application_id")

	req := &rpc.RemoteStopTransactionReqPublic{}
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

	res, err := api.operationService.RemoteStopTransaction(ctx, &rpc.RemoteStopTransactionReq{
		ApplicationId:         applicationId,
		ChargePointIdentifier: req.ChargePointIdentifier,
		ConnectorId:           req.ConnectorId,
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

// GetLatestStatus fetches the latest status notifications from each connector for a particular charge_point_id
// @Summary Get Status of Connectors
// @Tags Operations
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Body body rpc.GetLatestStatusNotificationsReqPublic true "Post GetLatestStatus body"
// @Success 200 {object} rpc.GetLatestStatusNotificationsResp
// @Router /v2/ocpp/operations/get-latest-status [post]
func (api *OperationsAPI) GetLatestStatus(c *gin.Context) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.CtxKey("gin"), c)

	applicationId := c.GetString("application_id")

	req := &rpc.GetLatestStatusNotificationsReqPublic{}
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

	res, err := api.statusnotificationService.GetLatestStatusNotifications(ctx, &rpc.GetLatestStatusNotificationsReq{
		ApplicationId:         applicationId,
		ChargePointIdentifier: req.ChargePointIdentifier,
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
