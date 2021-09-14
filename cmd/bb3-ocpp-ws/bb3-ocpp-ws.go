package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Beep-Technologies/beepbeep3-iam/pkg/logger"
	"github.com/Beep-Technologies/beepbeep3-ocpp/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	logger.LogInit()

	r := mux.NewRouter()
	r.HandleFunc("/{chargePointIdentifier}", handlers.HttpUpgradeHandler)

	host := "0.0.0.0"
	port := 8060
	addr := fmt.Sprintf("%s:%d", host, port)

	s := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	fmt.Printf("starting server at %s ...\n", addr)

	if err := s.ListenAndServe(); err != nil {
		fmt.Printf("%v\n", err)
	}
}
