package router

import (
	"os"

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
	operationsAPI       *controller.OperationsAPI
	applicationsAPI     *controller.ApplicationsAPI
	chargepointsAPI     *controller.ChargePointsAPI
}

func NewRouter(
	s *ocppserver.OCPPWebSocketServer,
	oa *controller.OperationsAPI,
	aa *controller.ApplicationsAPI,
	ca *controller.ChargePointsAPI,
) (rt *Router) {
	return &Router{
		ocppWebSocketServer: s,
		operationsAPI:       oa,
		applicationsAPI:     aa,
		chargepointsAPI:     ca,
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
	r.GET("/ocpp-central-system/:applicationId/:entityCode/:chargePointIdentifier", rt.ocppWebSocketServer.HttpUpgradeHandler)

	rg := r.Group("v2/ocpp")

	// for swagger
	hostUrl := os.Getenv("HOST_URL")
	if hostUrl == "" || hostUrl == "localhost:8060" {
		docs.SwaggerInfo.Host = "localhost:8060"
		docs.SwaggerInfo.Schemes = []string{"http"}
	} else {
		docs.SwaggerInfo.Host = hostUrl
		docs.SwaggerInfo.Schemes = []string{"https"}
	}

	rg.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// set up APIs
	rg.POST("/operations/remote-start-transaction", rt.operationsAPI.RemoteStartTransaction)
	rg.POST("/operations/remote-stop-transaction", rt.operationsAPI.RemoteStopTransaction)
	rg.POST("/operations/get-latest-status", rt.operationsAPI.GetLatestStatus)

	rg.POST("/charge_points", rt.chargepointsAPI.CreateChargePoint)
	rg.POST("/charge_points/id_tags", rt.chargepointsAPI.CreateChargePointIdTag)

	rg.POST("/applications", rt.applicationsAPI.CreateApplication)
	rg.POST("/applications/callbacks", rt.applicationsAPI.SetApplicationCallback)

	return r
}
