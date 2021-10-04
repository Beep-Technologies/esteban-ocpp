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
}

func NewRouter(s *ocppserver.OCPPWebSocketServer, oa *controller.OperationsAPI) (rt *Router) {
	return &Router{
		ocppWebSocketServer: s,
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
	r.GET("/ocpp/:chargePointIdentifier", rt.ocppWebSocketServer.HttpUpgradeHandler)

	rg := r.Group("v2/ocpp")

	// for swagger
	// set host to localhost if not live, else defaults to dev.beepbeep.tech
	if os.Getenv("LIVE") == "" {
		docs.SwaggerInfo.Host = "localhost:8060"
		docs.SwaggerInfo.Schemes = []string{"http"}
	}
	rg.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// set up APIs
	rg.POST("/operations/remote-start-transaction", rt.operationsAPI.RemoteStartTransaction)
	rg.POST("/operations/remote-stop-transaction", rt.operationsAPI.RemoteStopTransaction)
	rg.POST("/operations/get-latest-status", rt.operationsAPI.GetLatestStatus)

	return r
}
