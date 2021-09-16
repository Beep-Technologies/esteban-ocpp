package connection

import (
	"errors"
)

var ErrorFirstRequest = errors.New("first request sent from charge point must be a BootNotification request")
var ErrorOngoingOperation = errors.New("cannot set the operation as there is already an ongoing operation")
