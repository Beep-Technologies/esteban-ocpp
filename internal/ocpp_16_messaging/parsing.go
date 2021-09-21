package messaging

import (
	"encoding/json"
)

// GetOCPP16MessageType parses the WebSocket message body to get the message type
// GetOCPP16MessageType only checks if the message is a JSON array, and the MessageTypeId
func GetOCPP16MessageType(p []byte) (OCPP16MessageType, error) {
	// OCPP 1.6 messages are JSON arrays, with the first element
	// determining the message type

	arr := []interface{}{}
	err := json.Unmarshal(p, &arr)
	if err != nil {
		return 0, ErrorOCPP16Parse
	}

	if len(arr) < 3 {
		return 0, ErrorOCPP16Parse
	}

	messageTypeIDf, ok := arr[0].(float64)
	if !ok {
		return 0, ErrorOCPP16Parse
	}

	messageTypeID := OCPP16MessageType(messageTypeIDf)
	if !IsOCPP16MessageType(messageTypeID) {
		return 0, ErrorOCPP16Parse
	}

	return messageTypeID, nil
}

// ParseOCPP16Call parses the WebSocket message body
// ParseOCPP16Call does not validate if the values are correct, only the types
// ParseOCPP16Call also does not validate the payload body
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

// ParseOCPP16CallResult parses the WebSocket message body
// ParseOCPP16CallResult does not validate if the values are correct, only the types
// ParseOCPP16CallResult also does not validate the payload body
func ParseOCPP16CallResult(p []byte) (*OCPP16CallResult, error) {
	// an OCPP 1.6 call result message has the syntax
	// [<MessageTypeId>, "<UniqueId>", {<Payload>}]

	arr := []interface{}{}
	json.Unmarshal(p, &arr)

	if len(arr) != 3 {
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

	payload := arr[2]

	msg := &OCPP16CallResult{
		MessageTypeID: messageTypeID,
		UniqueID:      uniqueID,
		Payload:       payload,
	}

	return msg, nil
}

// ParseOCPP16CallError parses the WebSocket message body
// ParseOCPP16CallError does not validate if the values are correct, only the types
// ParseOCPP16CallError also does not validate the payload body
func ParseOCPP16CallError(p []byte) (*OCPP16CallError, error) {
	// an OCPP 1.6 call error message has the syntax
	// [<MessageTypeId>, "<UniqueId>", "<errorCode>", "<errorDescription>", {<errorDetails>}]

	arr := []interface{}{}
	json.Unmarshal(p, &arr)

	if len(arr) != 5 {
		return nil, ErrorOCPP16Parse
	}

	messageTypeID, ok := arr[0].(OCPP16MessageType)
	if !ok {
		return nil, ErrorOCPP16Parse
	}

	uniqueID, ok := arr[1].(string)
	if !ok {
		return nil, ErrorOCPP16Parse
	}

	errorCode, ok := arr[2].(OCPP16CallErrorCode)
	if !ok || !IsOCPP16CallErrorCode(errorCode) {
		return nil, ErrorOCPP16Parse
	}

	errorDescription, ok := arr[3].(string)
	if !ok {
		return nil, ErrorOCPP16Parse
	}

	errorDetails := arr[4]

	msg := &OCPP16CallError{
		MessageTypeID:    messageTypeID,
		UniqueID:         uniqueID,
		ErrorCode:        errorCode,
		ErrorDescription: errorDescription,
		ErrorDetails:     errorDetails,
	}

	return msg, nil
}

// ParseOCPP16 call unparses the struct and marshals it into valid JSON format
func UnparseOCPP16Call(c OCPP16CallMessage) ([]byte, error) {
	// an OCPP 1.6 call message has the syntax
	// [<MessageTypeId>, "<UniqueId>", "<Action>", {<Payload>}]

	var a [4]interface{}
	a[0] = c.MessageTypeID
	a[1] = c.UniqueID
	a[2] = c.Action
	a[3] = c.Payload

	p, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	return p, nil
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
