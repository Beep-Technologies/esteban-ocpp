package chargepoint

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging"
	application "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/application"
	chargepoint "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/charge_point"
	statusnotification "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/status_notification"
	transaction "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/transaction"
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

func (cp *OCPP16ChargePoint) RemoteStartTransaction() (transactionID int, err error) {
	// STUB
	return 0, nil
}

func (cp *OCPP16ChargePoint) RemoteStopTransaction() (err error) {
	// STUB
	return nil
}

func (cp *OCPP16ChargePoint) TriggerStatusNotification() (err error) {
	// STUB
	return nil
}