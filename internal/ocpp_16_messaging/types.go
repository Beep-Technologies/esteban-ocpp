package messaging

type OCPP16MessageType int

const (
	CALL       OCPP16MessageType = 2
	CALLRESULT OCPP16MessageType = 3
	CALLERROR  OCPP16MessageType = 4
)

func IsOCPP16MessageType(s OCPP16MessageType) bool {
	if s == CALL ||
		s == CALLRESULT ||
		s == CALLERROR {
		return true
	}

	return false
}

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

func IsOCPP16CallErrorCode(s OCPP16CallErrorCode) bool {
	if s == NotImplemented ||
		s == NotSupported ||
		s == InternalError ||
		s == ProtocolError ||
		s == SecurityError ||
		s == FormationViolation ||
		s == PropertyConstraintViolation ||
		s == OccurenceConstraintViolation ||
		s == TypeConstraintViolation ||
		s == GenericError {
		return true
	}

	return false
}

type OCPP16CallError struct {
	MessageTypeID    OCPP16MessageType
	UniqueID         string
	ErrorCode        OCPP16CallErrorCode
	ErrorDescription string
	ErrorDetails     interface{}
}

type OCPP16Status string

const (
	Available     OCPP16Status = "Available"
	Preparing     OCPP16Status = "Preparing"
	Charging      OCPP16Status = "Charging"
	SuspendedEV   OCPP16Status = "SuspendedEV"
	SuspendedEVSE OCPP16Status = "SuspendedEVSE"
	Finishing     OCPP16Status = "Finishing"
	Reserved      OCPP16Status = "Reserved"
	Unavailable   OCPP16Status = "Unavailable"
	Faulted       OCPP16Status = "Faulted"
)
