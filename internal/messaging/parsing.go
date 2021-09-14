package messaging

import (
	"encoding/json"
	"errors"
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
		return nil, errors.New("len(arr)!=4")
	}

	messageTypeIDf, ok := arr[0].(float64)
	if !ok {
		return nil, errors.New("messageTypeID cannot converted")
	}
	messageTypeID := OCPP16MessageType(messageTypeIDf)

	uniqueID, ok := arr[1].(string)
	if !ok {
		return nil, errors.New("uniqueID cannot converted")
	}

	action, ok := arr[2].(string)
	if !ok {
		return nil, errors.New("action cannot converted")
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
