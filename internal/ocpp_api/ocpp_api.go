package ocpp_api

import "log"

type OCPPWebSocketApp struct {
	logger *log.Logger
}

func NewOCPPWebSocketApp(l *log.Logger) *OCPPWebSocketApp {
	return &OCPPWebSocketApp{
		logger: l,
	}
}
