package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/Beep-Technologies/beepbeep3-iam/pkg/db"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rest/controller"
	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rest/router"

	ocpp16cs "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16_cs"
	ocppserver "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_server"
	chargepointsrv "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/charge_point"
	operationsrv "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/operation"
	statusnotificationsrv "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/status_notification"
	transactionsrv "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/transaction"
)

// @title BB3 OCPP API
// @version 2.0
// @description Service to interface with OCPP-compliant charge points

// @contact.name Lowen
// @contact.email lowen@beepbeep.tech

// @host dev.beepbeep.tech
// @BasePath /v2
// @Schemes https
func main() {
	// load from .env file. doesn't matter if an error is returned
	godotenv.Load()

	db.ConnectDataBase()

	l := log.New(os.Stdout, "", 0)

	chargePointService := chargepointsrv.NewService(db.ORM)
	transactionService := transactionsrv.NewService(db.ORM)
	statusNotificationService := statusnotificationsrv.NewService(db.ORM)

	ocpp16CentralSystem := ocpp16cs.NewOCPP16CentralSystem(
		l,
		chargePointService,
		transactionService,
		statusNotificationService,
	)

	operationService := operationsrv.NewService(
		db.ORM,
		chargePointService,
		transactionService,
		ocpp16CentralSystem,
	)

	ocppWebSocketServer := ocppserver.NewOCPPWebSocketServer(l, ocpp16CentralSystem)
	operationController := controller.NewOperationsAPI(operationService)

	rt := router.NewRouter(ocppWebSocketServer, operationController)

	r := gin.New()
	r.Use(gin.LoggerWithWriter(l.Writer()))
	rt.Apply(r)

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
