package messaging

import (
	"encoding/json"
)

// ParseOCPP16 call parses the WebSocket message body
// ParseOCPP16 does not validate if the values are correct, only the types
// ParseOCPP16 also does not validate the payload body
func ParseOCPP16Call(p []byte) (*OCPP16CallMessage, error) {
	// an OCPP 1.6 call message has the syntax
	// [<MessageTypeId>, "<UniqueId>", "<Action>", {<Payload>}]

	arr := []interface{}{}
	json.Unmarshal(p, &arr)

	if len(arr) != 4 {
		return nil, ErrorOCPP16Parse
	}

	messageTypeIDf, ok := arr[0].(float64)
	if !ok {
		return nil, ErrorOCPP16Parse
	}
	messageTypeID := OCPP16MessageType(messageTypeIDf)

	uniqueID, ok := arr[1].(string)
	if !ok {
		return nil, ErrorOCPP16Parse
	}

	action, ok := arr[2].(string)
	if !ok {
		return nil, ErrorOCPP16Parse
	}

	payload := arr[3]

	msg := &OCPP16CallMessage{
		MessageTypeID: messageTypeID,
		UniqueID:      uniqueID,
		Action:        action,
		Payload:       payload,
	}

	return msg, nil
}

// UnparseOCPP16CallResult unparses the struct and marshals it into valid JSON format
func UnparseOCPP16CallResult(r *OCPP16CallResult) ([]byte, error) {
	// an OCPP 1.6 call result message has the syntax
	// [<MessageTypeId>, "<UniqueId>", {<Payload>}]

	var a [3]interface{}
	a[0] = r.MessageTypeID
	a[1] = r.UniqueID
	a[2] = r.Payload

	p, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// UnparseOCPP16CallError unparses the struct and marshals it into valid JSON format
func UnparseOCPP16CallError(r *OCPP16CallError) ([]byte, error) {
	// an OCPP 1.6 call error message has the syntax
	// [<MessageTypeId>, "<UniqueId>", "<errorCode>", "<errorDescription>", {<errorDetails>}]

	var a [5]interface{}
	a[0] = r.MessageTypeID
	a[1] = r.UniqueID
	a[2] = r.ErrorCode
	a[3] = r.ErrorDescription
	a[4] = r.ErrorDetails

	p, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	return p, nil
}
