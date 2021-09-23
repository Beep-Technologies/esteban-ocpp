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

	o16cs := ocpp16cs.NewOCPP16CentralSystem(l)      // ocpp 1.6 central system
	o := ocppserver.NewOCPPWebSocketServer(l, o16cs) // ocpp websocket server
	oa := controller.NewOperationsAPI(o16cs)         // operations api controller
	rt := router.NewRouter(o, oa)                    // api router

	r := gin.Default()
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
