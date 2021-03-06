package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/operation"
	statusnotification "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/status_notification"
	"github.com/Beep-Technologies/beepbeep3-ocpp/pkg/constants"
	"github.com/Beep-Technologies/beepbeep3-ocpp/pkg/logger"
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
// @Param Body body rpc.RemoteStartTransactionReq true "Post RemoteStartTransactionReq body"
// @Success 200 {object} rpc.RemoteStartTransactionResp
// @Router /v2/ocpp/operations/remote-start-transaction [post]
func (api *OperationsAPI) RemoteStartTransaction(c *gin.Context) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.CtxKey("gin"), c)

	req := &rpc.RemoteStartTransactionReq{}
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

	res, err := api.operationService.RemoteStartTransaction(ctx, req)

	if err != nil {
		logger.Error(
			err.Error(),
			zap.String("event", "RemoteStartTransaction"),
			zap.String("error", err.Error()),
		)
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
// @Param Body body rpc.RemoteStopTransactionReq true "Post RemoteStopTransactionReq body"
// @Success 200 {object} rpc.RemoteStopTransactionResp
// @Router /v2/ocpp/operations/remote-stop-transaction [post]
func (api *OperationsAPI) RemoteStopTransaction(c *gin.Context) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.CtxKey("gin"), c)

	req := &rpc.RemoteStopTransactionReq{}
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

	res, err := api.operationService.RemoteStopTransaction(ctx, req)

	if err != nil {
		logger.Error(
			err.Error(),
			zap.String("event", "RemoteStopTransaction"),
			zap.String("error", err.Error()),
		)
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
