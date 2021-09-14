// Code generated by schema-generate. DO NOT EDIT.

package ocpp16

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

// StopTransactionRequest
type StopTransactionRequest struct {
	IdTag           string                  `json:"idTag,omitempty"`
	MeterStop       int                     `json:"meterStop"`
	Reason          string                  `json:"reason,omitempty"`
	Timestamp       string                  `json:"timestamp"`
	TransactionData []*TransactionDataItems `json:"transactionData,omitempty"`
	TransactionId   int                     `json:"transactionId"`
}

// TransactionDataItems
type TransactionDataItems struct {
	SampledValue []*SampledValueItems `json:"sampledValue"`
	Timestamp    string               `json:"timestamp"`
}

func (strct *StopTransactionRequest) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "idTag" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"idTag\": ")
	if tmp, err := json.Marshal(strct.IdTag); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// "MeterStop" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "meterStop" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"meterStop\": ")
	if tmp, err := json.Marshal(strct.MeterStop); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "reason" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"reason\": ")
	if tmp, err := json.Marshal(strct.Reason); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// "Timestamp" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "timestamp" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"timestamp\": ")
	if tmp, err := json.Marshal(strct.Timestamp); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "transactionData" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"transactionData\": ")
	if tmp, err := json.Marshal(strct.TransactionData); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// "TransactionId" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "transactionId" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"transactionId\": ")
	if tmp, err := json.Marshal(strct.TransactionId); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *StopTransactionRequest) UnmarshalJSON(b []byte) error {
	meterStopReceived := false
	timestampReceived := false
	transactionIdReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "idTag":
			if err := json.Unmarshal([]byte(v), &strct.IdTag); err != nil {
				return err
			}
		case "meterStop":
			if err := json.Unmarshal([]byte(v), &strct.MeterStop); err != nil {
				return err
			}
			meterStopReceived = true
		case "reason":
			if err := json.Unmarshal([]byte(v), &strct.Reason); err != nil {
				return err
			}
		case "timestamp":
			if err := json.Unmarshal([]byte(v), &strct.Timestamp); err != nil {
				return err
			}
			timestampReceived = true
		case "transactionData":
			if err := json.Unmarshal([]byte(v), &strct.TransactionData); err != nil {
				return err
			}
		case "transactionId":
			if err := json.Unmarshal([]byte(v), &strct.TransactionId); err != nil {
				return err
			}
			transactionIdReceived = true
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if meterStop (a required property) was received
	if !meterStopReceived {
		return errors.New("\"meterStop\" is required but was not present")
	}
	// check if timestamp (a required property) was received
	if !timestampReceived {
		return errors.New("\"timestamp\" is required but was not present")
	}
	// check if transactionId (a required property) was received
	if !transactionIdReceived {
		return errors.New("\"transactionId\" is required but was not present")
	}
	return nil
}

func (strct *TransactionDataItems) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// "SampledValue" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "sampledValue" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"sampledValue\": ")
	if tmp, err := json.Marshal(strct.SampledValue); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// "Timestamp" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "timestamp" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"timestamp\": ")
	if tmp, err := json.Marshal(strct.Timestamp); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *TransactionDataItems) UnmarshalJSON(b []byte) error {
	sampledValueReceived := false
	timestampReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "sampledValue":
			if err := json.Unmarshal([]byte(v), &strct.SampledValue); err != nil {
				return err
			}
			sampledValueReceived = true
		case "timestamp":
			if err := json.Unmarshal([]byte(v), &strct.Timestamp); err != nil {
				return err
			}
			timestampReceived = true
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if sampledValue (a required property) was received
	if !sampledValueReceived {
		return errors.New("\"sampledValue\" is required but was not present")
	}
	// check if timestamp (a required property) was received
	if !timestampReceived {
		return errors.New("\"timestamp\" is required but was not present")
	}
	return nil
}