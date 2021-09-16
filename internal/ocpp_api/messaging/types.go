package messaging

type OCPP16MessageType int

const (
	CALL       OCPP16MessageType = 2
	CALLRESULT OCPP16MessageType = 3
	CALLERROR  OCPP16MessageType = 4
)

type OCPP16CallMessage struct {
	MessageTypeID OCPP16MessageType
	UniqueID      string
	Action        string
	Payload       interface{}
}

type OCPP16CallResult struct {
	MessageTypeID OCPP16MessageType
	UniqueID      string
	Payload       interface{}
}

type OCPP16CallErrorCode string

const (
	NotImplemented               OCPP16CallErrorCode = "NotImplemented"
	NotSupported                 OCPP16CallErrorCode = "NotSupported"
	InternalError                OCPP16CallErrorCode = "InternalError"
	ProtocolError                OCPP16CallErrorCode = "ProtocolError"
	SecurityError                OCPP16CallErrorCode = "SecurityError"
	FormationViolation           OCPP16CallErrorCode = "FormationViolation"
	PropertyConstraintViolation  OCPP16CallErrorCode = "PropertyConstraintViolation"
	OccurenceConstraintViolation OCPP16CallErrorCode = "OccurenceConstraintViolation"
	TypeConstraintViolation      OCPP16CallErrorCode = "TypeConstraintViolation"
	GenericError                 OCPP16CallErrorCode = "GenericError"
)

type OCPP16CallError struct {
	MessageTypeID    OCPP16MessageType
	UniqueID         string
	ErrorCode        OCPP16CallErrorCode
	ErrorDescription string
	ErrorDetails     interface{}
}
