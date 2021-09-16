package operations

type OCPP16OperationType int

const (
	// these two operations are not CP / CS intiated
	NONE OCPP16OperationType = iota
	INITIAL
	DATA_TRANSFER
	BOOT_NOTIFICATION
	DIAGNOSTICS_STATUS_NOTIFICATION
	FIRMWARE_STATUS_NOTIFICATION
	HEARTBEAT
	METER_VALUES
	START_TRANSACTION
	STATUS_NOTIFICATION
	STOP_TRANSACTION
	CANCEL_RESERVATION
	CHANGE_AVAILABILITY
	CHANGE_CONFIGURATION
	CLEAR_CACHE
	CLEAR_CHARGING_PROFILE
	GET_COMPOSITE_SCHEDULE
	GET_CONFIGURATION
	GET_DIAGNOSTICS
	GET_LOCAL_LIST_VERSION
	REMOTE_START_TRANSACTION
	REMOTE_STOP_TRANSACTION
	RESERVE_NOW
	RESET
	SEND_LOCAL_LIST
	SET_CHARGING_PROFILE
	TRIGGER_MESSAGE
	UNLOCK_CONNECTOR
	UPDATE_FIRMWARE
)

var chargePointInitiated map[OCPP16OperationType]bool = map[OCPP16OperationType]bool{
	// neither CP / CS intiated
	NONE:    false,
	INITIAL: false,

	// both CP & CS initiated
	DATA_TRANSFER: true,

	// charge point (CP) initiated
	BOOT_NOTIFICATION:               true,
	DIAGNOSTICS_STATUS_NOTIFICATION: true,
	FIRMWARE_STATUS_NOTIFICATION:    true,
	HEARTBEAT:                       true,
	METER_VALUES:                    true,
	START_TRANSACTION:               true,
	STATUS_NOTIFICATION:             true,
	STOP_TRANSACTION:                true,

	// central system (CS) initiated
	CANCEL_RESERVATION:       false,
	CHANGE_AVAILABILITY:      false,
	CHANGE_CONFIGURATION:     false,
	CLEAR_CACHE:              false,
	CLEAR_CHARGING_PROFILE:   false,
	GET_COMPOSITE_SCHEDULE:   false,
	GET_CONFIGURATION:        false,
	GET_DIAGNOSTICS:          false,
	GET_LOCAL_LIST_VERSION:   false,
	REMOTE_START_TRANSACTION: false,
	REMOTE_STOP_TRANSACTION:  false,
	RESERVE_NOW:              false,
	RESET:                    false,
	SEND_LOCAL_LIST:          false,
	SET_CHARGING_PROFILE:     false,
	TRIGGER_MESSAGE:          false,
	UNLOCK_CONNECTOR:         false,
	UPDATE_FIRMWARE:          false,
}

var centralSystemInitiated map[OCPP16OperationType]bool = map[OCPP16OperationType]bool{
	// neither CP / CS intiated
	NONE:    false,
	INITIAL: false,

	// both CP & CS initiated
	DATA_TRANSFER: true,

	// charge point (CP) initiated
	BOOT_NOTIFICATION:               false,
	DIAGNOSTICS_STATUS_NOTIFICATION: false,
	FIRMWARE_STATUS_NOTIFICATION:    false,
	HEARTBEAT:                       false,
	METER_VALUES:                    false,
	START_TRANSACTION:               false,
	STATUS_NOTIFICATION:             false,
	STOP_TRANSACTION:                false,

	// central system (CS) initiated
	CANCEL_RESERVATION:       true,
	CHANGE_AVAILABILITY:      true,
	CHANGE_CONFIGURATION:     true,
	CLEAR_CACHE:              true,
	CLEAR_CHARGING_PROFILE:   true,
	GET_COMPOSITE_SCHEDULE:   true,
	GET_CONFIGURATION:        true,
	GET_DIAGNOSTICS:          true,
	GET_LOCAL_LIST_VERSION:   true,
	REMOTE_START_TRANSACTION: true,
	REMOTE_STOP_TRANSACTION:  true,
	RESERVE_NOW:              true,
	RESET:                    true,
	SEND_LOCAL_LIST:          true,
	SET_CHARGING_PROFILE:     true,
	TRIGGER_MESSAGE:          true,
	UNLOCK_CONNECTOR:         true,
	UPDATE_FIRMWARE:          true,
}

func IsChargePointInitated(o OCPP16OperationType) bool {
	return chargePointInitiated[o]
}

func IsCentralSystemInitiated(o OCPP16OperationType) bool {
	return centralSystemInitiated[o]
}
