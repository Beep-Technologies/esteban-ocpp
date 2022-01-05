package router

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/Beep-Technologies/beepbeep3-ocpp/api/rest/controller"
	docs "github.com/Beep-Technologies/beepbeep3-ocpp/docs"
	middleware "github.com/Beep-Technologies/beepbeep3-ocpp/internal/middleware"
	ocppserver "github.com/Beep-Technologies/beepbeep3-ocpp/internal/ocpp_server"
)

type Router struct {
	ocppWebSocketServer *ocppserver.OCPPWebSocketServer
	applicationsAPI     *controller.ApplicationsAPI
	chargepointsAPI     *controller.ChargePointsAPI
	middleware          *middleware.Middleware
}

func NewRouter(
	s *ocppserver.OCPPWebSocketServer,
	aa *controller.ApplicationsAPI,
	ca *controller.ChargePointsAPI,
	m *middleware.Middleware,
) (rt *Router) {
	return &Router{
		ocppWebSocketServer: s,
		applicationsAPI:     aa,
		chargepointsAPI:     ca,
		middleware:          m,
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
	r.GET("/ocpp-central-system//:entityCode/:chargePointIdentifier", rt.ocppWebSocketServer.HttpUpgradeHandler)

	rg := r.Group("v2/ocpp")

	// for swagger
	hostUrl := os.Getenv("HOST_URL")
	if hostUrl == "" || hostUrl == "localhost:8060" {
		docs.SwaggerInfo.Host = "localhost:8060"
		docs.SwaggerInfo.Schemes = []string{"http"}
	} else {
		docs.SwaggerInfo.Host = hostUrl
		docs.SwaggerInfo.Schemes = []string{"https", "http"}
	}

	rg.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	rg.POST("/charge_points", rt.middleware.APIKeyMiddleware(), rt.chargepointsAPI.CreateChargePoint)
	rg.POST("/charge_points/id_tags", rt.middleware.APIKeyMiddleware(), rt.chargepointsAPI.CreateChargePointIdTag)

	rg.GET("/applications/callbacks", rt.middleware.APIKeyMiddleware(), rt.applicationsAPI.GetApplicationCallbacks)
	rg.POST("/applications/callbacks", rt.middleware.APIKeyMiddleware(), rt.applicationsAPI.SetApplicationCallback)

	return r
}
