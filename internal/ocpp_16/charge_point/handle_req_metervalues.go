package chargepoint

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rpc"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/messaging/schemas"
	"github.com/mitchellh/mapstructure"
)

type ByDate []*schemas.MeterValueItems

func (a ByDate) Len() int      { return len(a) }
func (a ByDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool {
	RFC3339Milli := "2006-01-02T15:04:05.000Z07:00"
	ait, aie := time.Parse(RFC3339Milli, a[i].Timestamp)
	ajt, aje := time.Parse(RFC3339Milli, a[j].Timestamp)

	if aie != nil || aje != nil {
		return false
	}

	return ait.Unix() < ajt.Unix()
}

func (cp *OCPP16ChargePoint) handleMeterValues(msg messaging.OCPP16CallMessage) (*messaging.OCPP16CallResult, *messaging.OCPP16CallError) {
	p := &schemas.MeterValuesRequest{}

	// decode the payload into the struct
	err := mapstructure.Decode(msg.Payload, p)
	if err != nil {
		return nil, &messaging.OCPP16CallError{
			MessageTypeID:    messaging.CALLERROR,
			UniqueID:         msg.UniqueID,
			ErrorCode:        messaging.FormationViolation,
			ErrorDescription: "",
			ErrorDetails:     struct{}{},
		}
	}

	// look through the meter values, make a callback if
	// there is a value for energyActiveImportRegister,
	// i.e. total energy supplied during charging session

	// sort meter values by timestamp
	// sort in reverse order i.e. latest meter values first
	sort.Sort(sort.Reverse(ByDate(p.MeterValue)))

	energyActiveImportRegister := 0
	energyActiveImportRegisterFound := false

	// get the latest value of energyActiveImportRegister
	// look through meter values, starting with the latest values
	for _, mv := range p.MeterValue {
		// look through the sampled value
		for _, sv := range mv.SampledValue {
			// default value of mv.Measurand is "Energy.Active.Import.Register"
			if sv.Measurand == "" ||
				sv.Measurand == "Energy.Active.Import.Register" {
				v, err := strconv.Atoi(sv.Value)
				if err != nil {
					continue
				}

				energyActiveImportRegister = v
				energyActiveImportRegisterFound = true
			}

			if energyActiveImportRegisterFound {
				break
			}
		}

		if energyActiveImportRegisterFound {
			break
		}
	}

	// if there is a value for energyActiveImportRegister, make the callback
	if energyActiveImportRegisterFound {
		transaction, err := cp.transactionService.GetOngoingTransaction(cp.ctx, &rpc.GetOngoingTransactionReq{
			EntityCode:            cp.entityCode,
			ChargePointIdentifier: cp.chargePointIdentifier,
			ConnectorId:           int32(p.ConnectorId),
		})

		if err != nil {
			cp.logger.Error(fmt.Sprintf("handleMeterValues for non existent transaction id %v", p.TransactionId))
		}

		// energyActiveImportRegister should be the DIFFERENCE in start_meter_value and actual energyActiveImportRegister (the current meter reading value)
		energyActiveImportRegister = energyActiveImportRegister - int(transaction.Transaction.StartMeterValue)

		d := map[string]interface{}{
			"meter_values": map[string]interface{}{
				"energy_active_import_register": energyActiveImportRegister,
			},
		}

		go cp.makeCallback("MeterValues", d)
	}

	return &messaging.OCPP16CallResult{
		MessageTypeID: messaging.CALLRESULT,
		UniqueID:      msg.UniqueID,
		Payload:       &schemas.MeterValuesResponse{},
	}, nil
}
