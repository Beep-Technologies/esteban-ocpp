package router

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rest/controller"
	docs "github.com/Beep-Technologies/beepbeep3-ocpp/docs"
	ocppserver "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_server"
)

type Router struct {
	ocppWebSocketServer *ocppserver.OCPPWebSocketServer
	chargepointsAPI     *controller.ChargePointsAPI
	operationsAPI       *controller.OperationsAPI
}

func NewRouter(
	s *ocppserver.OCPPWebSocketServer,
	ca *controller.ChargePointsAPI,
	oa *controller.OperationsAPI,
) (rt *Router) {
	return &Router{
		ocppWebSocketServer: s,
		chargepointsAPI:     ca,
		operationsAPI:       oa,
	}
}

func (rt *Router) Apply(r *gin.Engine) *gin.Engine {
	// CORS
	r.Use(cors.New(cors.Config{
		AllowCredentials: false,
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", " Content-Length", "Authorization", "accept", "origin", "Referer", "User-Agent"},
	}))

	// set up websocket server endpoint
	r.GET("/ocpp-central-system/:entityCode/:chargePointIdentifier", rt.ocppWebSocketServer.HttpUpgradeHandler)

	rg := r.Group("v2/ocpp")

	// for swagger
	hostUrl := os.Getenv("HOST_URL")
	port := 8060
	if hostPortStr := os.Getenv("HOST_PORT"); hostPortStr != "" {
		port, _ = strconv.Atoi(hostPortStr)
	}

	if hostUrl == "" || hostUrl == "localhost" {
		docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", hostUrl, port)
		docs.SwaggerInfo.Schemes = []string{"http"}
	} else {
		docs.SwaggerInfo.Host = hostUrl
		docs.SwaggerInfo.Schemes = []string{"https", "http"}
	}

	rg.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// set up APIs
	rg.POST("/operations/remote-start-transaction", rt.operationsAPI.RemoteStartTransaction)
	rg.POST("/operations/remote-stop-transaction", rt.operationsAPI.RemoteStopTransaction)

	rg.GET("/charge-points/:entityCode/:chargePointIdentifier", rt.chargepointsAPI.GetChargePoint)
	rg.GET("/charge-points/:entityCode/:chargePointIdentifier/status", rt.chargepointsAPI.GetChargePointConnectorStatus)

	return r
}
