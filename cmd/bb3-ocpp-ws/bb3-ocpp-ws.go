package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Beep-Technologies/beepbeep3-iam/pkg/logger"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_api"
	"github.com/gorilla/mux"
)

func main() {
	logger.LogInit()
	l := log.New(os.Stdout, "", (log.Ltime | log.Lmicroseconds))
	o := ocpp_api.NewOCPPWebSocketApp(l)

	r := mux.NewRouter()
	r.HandleFunc("/{chargePointIdentifier}", o.HttpUpgradeHandler)

	host := "0.0.0.0"
	port := 8060
	addr := fmt.Sprintf("%s:%d", host, port)

	s := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	l.Printf("starting server at %s ...\n", addr)

	if err := s.ListenAndServe(); err != nil {
		l.Printf("%v\n", err)
	}
}
