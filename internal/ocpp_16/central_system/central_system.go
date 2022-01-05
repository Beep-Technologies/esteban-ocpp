package centralsystem

import (
	"context"
	"errors"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	ocpp16 "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16"
	ocpp16cp "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/charge_point"
	messaging "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging"

	application "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/application"
	chargepoint "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/charge_point"
	statusnotification "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/status_notification"
	transaction "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/transaction"
)

type OCPP16CentralSystem struct {
	logger                    *zap.Logger
	chargePoints              map[string]ocpp16.ChargePoint // key is {entityCode}/{chargePointIdentifier}
	applicationService        *application.Service
	chargepointService        *chargepoint.Service
	transactionService        *transaction.Service
	statusNotificationService *statusnotification.Service
}

func NewOCPP16CentralSystem(
	logger *zap.Logger,
	applicationService *application.Service,
	chargepointService *chargepoint.Service,
	transactionService *transaction.Service,
	statusNotificationService *statusnotification.Service,
) *OCPP16CentralSystem {
	centralSystemLogger := logger.With(
		zap.String("source", "central_system"),
	)

	return &OCPP16CentralSystem{
		logger:                    centralSystemLogger,
		chargePoints:              make(map[string]ocpp16.ChargePoint),
		applicationService:        applicationService,
		chargepointService:        chargepointService,
		transactionService:        transactionService,
		statusNotificationService: statusNotificationService,
	}
}

// GetChargePoint gets a charge point that is currently connected to the central system
func (cs *OCPP16CentralSystem) GetChargePoint(entityCode, identifier string) (ocpp16.ChargePoint, error) {
	key := entityCode + "/" + identifier
	cp, ok := cs.chargePoints[key]
	if !ok {
		return nil, errors.New("charge point not found")
	}

	return cp, nil
}

// ConnectChargePoint returns an error if there is an error during the initial connection to the charge point,
// else it returns with nil after the charge point's context is cancelled
func (cs *OCPP16CentralSystem) ConnectChargePoint(entityCode, identifier string, conn *websocket.Conn) error {
	key := entityCode + "/" + identifier
	// check if the key already exists
	_, ok := cs.chargePoints[key]
	if ok {
		return errors.New("there is already a charge point connected with the key " + key)
	}

	// check if charge point exists
	cpRes, err := cs.chargepointService.GetChargePoint(context.Background(), &rpc.GetChargePointReq{
		EntityCode:            entityCode,
		ChargePointIdentifier: identifier,
	})
	if err != nil {
		return err
	}

	// if charge point exists, attach the charge point to this central system
	ctx, ctxCancel := context.WithCancel(context.Background())
	inCallStream := make(chan messaging.OCPP16CallMessage)
	inCallResultStream := make(chan messaging.OCPP16CallResult)
	inCallErrorStream := make(chan messaging.OCPP16CallError)
	outCallStream := make(chan messaging.OCPP16CallMessage)
	outCallResultStream := make(chan messaging.OCPP16CallResult)
	outCallErrorStream := make(chan messaging.OCPP16CallError)

	cp := ocpp16cp.NewOCPP16ChargePoint(
		int(cpRes.ChargePoint.Id),
		identifier,
		entityCode,
		ctx,
		inCallStream,
		inCallResultStream,
		inCallErrorStream,
		outCallStream,
		outCallResultStream,
		outCallErrorStream,
		cs.logger,
		cs.applicationService,
		cs.chargepointService,
		cs.transactionService,
		cs.statusNotificationService,
	)

	cs.chargePoints[key] = cp
	go cp.Listen()

	readPump := func() {
		defer ctxCancel()

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			_, p, err := conn.ReadMessage()
			if err != nil {
				cs.logger.Error(
					err.Error(),
					zap.String("event", "receive_message_error"),
					zap.String("charge_point_key", key),
				)
				return
			}

			cs.logger.Info(
				string(p),
				zap.String("event", "receive_message"),
				zap.String("charge_point_key", key),
			)

			msgType, err := messaging.GetOCPP16MessageType(p)
			if err != nil {
				cs.logger.Error(
					err.Error(),
					zap.String("event", "parse_message_error"),
					zap.String("charge_point_key", key),
				)
				return
			}

			switch msgType {
			case messaging.CALL:
				msg, err := messaging.ParseOCPP16Call(p)
				if err != nil {
					cs.logger.Error(
						err.Error(),
						zap.String("event", "parse_call_message_error"),
						zap.String("charge_point_key", key),
					)
					return
				}
				inCallStream <- *msg
			case messaging.CALLRESULT:
				msg, err := messaging.ParseOCPP16CallResult(p)
				if err != nil {
					cs.logger.Error(
						err.Error(),
						zap.String("event", "parse_call_result_message_error"),
						zap.String("charge_point_key", key),
					)
					return
				}
				inCallResultStream <- *msg
			case messaging.CALLERROR:
				msg, err := messaging.ParseOCPP16CallError(p)
				if err != nil {
					cs.logger.Error(
						err.Error(),
						zap.String("event", "parse_call_error_message_error"),
						zap.String("charge_point_key", key),
					)
					return
				}
				inCallErrorStream <- *msg
			}
		}
	}

	writePump := func() {
		defer ctxCancel()
		for {
			var res []byte
			var err error
			select {
			case <-ctx.Done():
				return
			case msg := <-outCallStream:
				res, err = messaging.UnparseOCPP16Call(msg)
				if err != nil {
					cs.logger.Error(
						err.Error(),
						zap.String("event", "unparse_call_message_error"),
						zap.String("charge_point_key", key),
					)
					return
				}
			case msg := <-outCallResultStream:
				res, err = messaging.UnparseOCPP16CallResult(msg)
				if err != nil {
					cs.logger.Error(
						err.Error(),
						zap.String("event", "unparse_call_result_message_error"),
						zap.String("charge_point_key", key),
					)
					return
				}
			case msg := <-outCallErrorStream:
				res, err = messaging.UnparseOCPP16CallError(msg)
				if err != nil {
					cs.logger.Error(
						err.Error(),
						zap.String("event", "unparse_call_error_message_error"),
						zap.String("charge_point_key", key),
					)
					return
				}
			}

			err = conn.WriteMessage(websocket.TextMessage, res)
			cs.logger.Info(
				string(res),
				zap.String("event", "send_message"),
				zap.String("charge_point_key", key),
			)
			if err != nil {
				cs.logger.Error(
					err.Error(),
					zap.String("event", "websocket_write_message_error"),
					zap.String("charge_point_key", key),
				)
				return
			}
		}
	}

	go readPump()
	go writePump()

	// wait for context to be closed
	// this happens if the websocket connection is closed by the charge point
	// or if there is some error
	<-ctx.Done()

	// cleanup
	delete(cs.chargePoints, key)
	conn.Close()
	cs.logger.Info(
		"",
		zap.String("event", "websocket_closed"),
		zap.String("charge_point_key", key),
	)

	return nil
}
