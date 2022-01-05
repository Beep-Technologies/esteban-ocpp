package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/Beep-Technologies/beepbeep3-iam/pkg/db"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rest/controller"
	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rest/router"
	middleware "github.com/Beep-Technologies/beepbeep3-ocpp/internal/middleware"
	ocpp16cs "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_16/central_system"
	ocppserver "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_server"
	applicationsrv "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/application"
	chargepointsrv "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/charge_point"
	statusnotificationsrv "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/status_notification"
	transactionsrv "github.com/Beep-Technologies/beepbeep3-ocpp/internal/service/transaction"
)

// @title BB3 OCPP API
// @version 2.0
// @description Service to interface with OCPP-compliant charge points

// @contact.name Lowen
// @contact.email lowen@beepbeep.tech

// @host ocpp-dev.chargegowhere.sg:8060
// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// load from .env file. doesn't matter if an error is returned
	godotenv.Load()

	db.ConnectDataBase()

	var l *zap.Logger
	if os.Getenv("ENVIRONMENT") == "production" {
		os.Setenv("GIN_MODE", "release")
		l, _ = zap.NewProduction()
	} else {
		l, _ = zap.NewDevelopment()
	}
	defer l.Sync()

	applicationService := applicationsrv.NewService(db.ORM)
	chargePointService := chargepointsrv.NewService(db.ORM)
	transactionService := transactionsrv.NewService(db.ORM)
	statusNotificationService := statusnotificationsrv.NewService(db.ORM)

	ocpp16CentralSystem := ocpp16cs.NewOCPP16CentralSystem(
		l,
		applicationService,
		chargePointService,
		transactionService,
		statusNotificationService,
	)

	// operationService := operationsrv.NewService(
	// 	db.ORM,
	// 	chargePointService,
	// 	transactionService,
	// 	ocpp16CentralSystem,
	// )

	ocppWebSocketServer := ocppserver.NewOCPPWebSocketServer(l, ocpp16CentralSystem)
	// operationController := controller.NewOperationsAPI(
	// 	operationService,
	// 	statusNotificationService,
	// )
	applicationController := controller.NewApplicationsAPI(applicationService)
	chargepointController := controller.NewChargePointsAPI(chargePointService)

	rt := router.NewRouter(
		ocppWebSocketServer,
		applicationController,
		chargepointController,
		middleware.NewMiddleware(db.ORM),
	)

	// set up gin
	r := gin.New()
	// use ginzap middleware
	gl := l.With(zap.String("source", "gin"))
	r.Use(ginzap.Ginzap(gl, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(l, true))
	// apply router
	rt.Apply(r)

	host := "0.0.0.0"
	port := 80
	addr := fmt.Sprintf("%s:%d", host, port)

	s := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	l.Info("Starting server at " + addr + "...")
	if err := s.ListenAndServe(); err != nil {
		l.Error("Server error", zap.String("error", err.Error()))
	}
}
