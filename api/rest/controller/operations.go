package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Beep-Technologies/beepbeep3-iam/pkg/constants"
	"github.com/Beep-Technologies/beepbeep3-iam/pkg/db"
	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	ocpp16cs "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16_cs"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/operations"
)

type OperationsAPI struct {
	ocpp16CentralSystem *ocpp16cs.OCPP16CentralSystem
}

func NewOperationsAPI(o16cs *ocpp16cs.OCPP16CentralSystem) *OperationsAPI {
	return &OperationsAPI{
		ocpp16CentralSystem: o16cs,
	}
}

// RemoteStartTransaction has the Central System request the Charge Point to start a transaction
// @Summary Request Charge Point to Start a Transaction
// @Tags Operations
// @Accept json
// @Produce json
// @Param Body body rpc.RemoteStartTransactionReq true "Post RemoteStartTransactionReq body"
// @Success 200 {object} rpc.RemoteStartTransactionResp
// @Router /ocpp/operations/remote-start-transaction [post]
func (api *OperationsAPI) RemoteStartTransaction(c *gin.Context) {
	var req = rpc.RemoteStartTransactionReq{}
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

	srv := operations.NewService(db.ORM, api.ocpp16CentralSystem)
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.CtxKey("gin"), c)
	res, err := srv.RemoteStartTransaction(ctx, req)

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
// @Param Body body rpc.RemoteStopTransactionReq true "Post RemoteStopTransactionReq body"
// @Success 200 {object} rpc.RemoteStopTransactionResp
// @Router /ocpp/operations/remote-stop-transaction [post]
func (api *OperationsAPI) RemoteStopTransaction(c *gin.Context) {
	var req = rpc.RemoteStopTransactionReq{}
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

	srv := operations.NewService(db.ORM, api.ocpp16CentralSystem)
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.CtxKey("gin"), c)
	res, err := srv.RemoteStopTransaction(ctx, req)

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
