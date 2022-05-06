package chargepoint

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging/schemas"
	application "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/application"
	chargepoint "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/charge_point"
	statusnotification "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/status_notification"
	transaction "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/transaction"
	"github.com/google/uuid"
)

type OCPP16ChargePoint struct {
	// charge point metadata
	id                    int
	chargePointIdentifier string
	entityCode            string
	// context (mostly used for ctx.Done())
	ctx context.Context
	// input / output channels
	inCallStream        <-chan messaging.OCPP16CallMessage
	inCallResultStream  <-chan messaging.OCPP16CallResult
	inCallErrorStream   <-chan messaging.OCPP16CallError
	outCallStream       chan<- messaging.OCPP16CallMessage
	outCallResultStream chan<- messaging.OCPP16CallResult
	outCallErrorStream  chan<- messaging.OCPP16CallError
	// logger
	logger *zap.Logger
	// services
	applicationService        *application.Service
	chargepointService        *chargepoint.Service
	transactionService        *transaction.Service
	statusNotificationService *statusnotification.Service
	// for making OCPP RPC calls. currentCall refers to the current outgoing call.
	callMessageQueue *callMessageQueue
}

func NewOCPP16ChargePoint(
	id int,
	chargePointIdentifier string,
	entityCode string,
	ctx context.Context,
	inCallStream <-chan messaging.OCPP16CallMessage,
	inCallResultStream <-chan messaging.OCPP16CallResult,
	inCallErrorStream <-chan messaging.OCPP16CallError,
	outCallStream chan<- messaging.OCPP16CallMessage,
	outCallResultStream chan<- messaging.OCPP16CallResult,
	outCallErrorStream chan<- messaging.OCPP16CallError,
	logger *zap.Logger,
	applicationService *application.Service,
	chargepointService *chargepoint.Service,
	transactionService *transaction.Service,
	statusNotificationService *statusnotification.Service,
) *OCPP16ChargePoint {
	chargePointLogger := logger.With(
		zap.String("source", "charge_point"),
		zap.String("entity_code", entityCode),
		zap.String("charge_point_identifier", chargePointIdentifier),
		zap.String("charge_point_key", entityCode+"/"+chargePointIdentifier),
	)

	return &OCPP16ChargePoint{
		id:                        id,
		chargePointIdentifier:     chargePointIdentifier,
		entityCode:                entityCode,
		ctx:                       ctx,
		inCallStream:              inCallStream,
		inCallResultStream:        inCallResultStream,
		inCallErrorStream:         inCallErrorStream,
		outCallStream:             outCallStream,
		outCallResultStream:       outCallResultStream,
		outCallErrorStream:        outCallErrorStream,
		logger:                    chargePointLogger,
		applicationService:        applicationService,
		chargepointService:        chargepointService,
		transactionService:        transactionService,
		statusNotificationService: statusNotificationService,
		callMessageQueue: &callMessageQueue{
			queue: make([]messaging.OCPP16CallMessage, 0),
			mutex: &sync.Mutex{},
			cond:  sync.NewCond(&sync.Mutex{}),
		},
	}
}

func (cp *OCPP16ChargePoint) Listen() {
	// listen on channels for CS-initiated operations
	go cp.listenCS()
	// listen on channels for CP-initiated operations
	go cp.listenCP()
}

func (cp *OCPP16ChargePoint) RemoteStartTransaction(connectorID int) (transactionID int, err error) {
	// check if there is a currently ongoing transaction
	currentTransactionRes, err := cp.transactionService.GetOngoingTransaction(cp.ctx, &rpc.GetOngoingTransactionReq{
		EntityCode:            cp.entityCode,
		ChargePointIdentifier: cp.chargePointIdentifier,
		ConnectorId:           int32(connectorID),
	})
	if err != nil {
		return 0, err
	}
	if currentTransactionRes.OngoingTransaction {
		return 0, errors.New("there is already an ongoing transaction")
	}

	// generate a unique id tag for this connection
	// id tags must be 20 characters long
	idTag := uuid.NewString()[0:20]

	// create the transaction
	transactionRes, err := cp.transactionService.CreateTransaction(cp.ctx, &rpc.CreateTransactionReq{
		EntityCode:            cp.entityCode,
		ChargePointIdentifier: cp.chargePointIdentifier,
		ConnectorId:           int32(connectorID),
		RemoteInitiated:       true,
		IdTag:                 idTag,
	})
	if err != nil {
		return 0, err
	}

	cp.logger.Info(
		"remote transaction created",
		zap.String("event", "remote_transaction_created"),
		zap.Int32("transaction_id", transactionRes.Transaction.Id),
	)

	// add the message to the queue
	cp.callMessageQueue.enqueue(messaging.OCPP16CallMessage{
		MessageTypeID: messaging.CALL,
		UniqueID:      uuid.NewString(),
		Action:        "RemoteStartTransaction",
		Payload: schemas.RemoteStartTransactionRequest{
			ConnectorId: int(transactionRes.Transaction.ConnectorId),
			IdTag:       transactionRes.Transaction.IdTag,
		},
	})

	return int(transactionRes.Transaction.Id), nil
}

func (cp *OCPP16ChargePoint) RemoteStopTransaction(connectorID int) (err error) {
	// check if there is a currently ongoing transaction
	currentTransactionRes, err := cp.transactionService.GetOngoingTransaction(cp.ctx, &rpc.GetOngoingTransactionReq{
		EntityCode:            cp.entityCode,
		ChargePointIdentifier: cp.chargePointIdentifier,
		ConnectorId:           int32(connectorID),
	})
	if err != nil {
		return err
	}
	if !currentTransactionRes.OngoingTransaction {
		return errors.New("there is no ongoing transaction")
	}

	// add the message to the queue
	cp.callMessageQueue.enqueue(messaging.OCPP16CallMessage{
		MessageTypeID: messaging.CALL,
		UniqueID:      uuid.NewString(),
		Action:        "RemoteStopTransaction",
		Payload: schemas.RemoteStopTransactionRequest{
			TransactionId: int(currentTransactionRes.Transaction.Id),
		},
	})

	return nil
}

func (cp *OCPP16ChargePoint) TriggerStatusNotification(connectorID int, errorStatusCode messaging.OCPP16ChargePointErrorCode) (err error) {
	p := &schemas.StatusNotificationRequest{
		ConnectorId:     connectorID,
		ErrorCode:       string(errorStatusCode),
		Info:            "",
		Status:          string(messaging.Unavailable),
		Timestamp:       time.Now().Format(time.RFC3339),
		VendorErrorCode: "",
		VendorId:        "",
	}
	sn, err := cp.statusNotificationService.CreateStatusNotification(
		cp.ctx,
		&rpc.CreateStatusNotificationReq{
			EntityCode:            cp.entityCode,
			ChargePointIdentifier: cp.chargePointIdentifier,
			ConnectorId:           int32(p.ConnectorId),
			ErrorCode:             p.ErrorCode,
			Info:                  p.Info,
			Status:                p.Status,
			Timestamp:             p.Timestamp,
			VendorId:              p.VendorId,
			VendorErrorCode:       p.VendorErrorCode,
		},
	)

	// make callback
	go cp.makeCallback("StatusNotification", sn)
	return nil
}
