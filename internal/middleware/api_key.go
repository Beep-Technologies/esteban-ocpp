package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"

	application "github.com/Beep-Technologies/beepbeep3-ocpp/internal/repository/application"
)

type Middleware struct {
	application application.BaseRepo
}

func NewMiddleware(db *gorm.DB) *Middleware {
	return &Middleware{
		application: application.NewBaseRepo(db),
	}
}

func (m Middleware) APIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("authorization")

		if authHeader == "" {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"status":  http.StatusUnauthorized,
					"message": "EmptyAuthHeader",
					"data":    empty.Empty{},
				},
			)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Token") {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"status":  http.StatusUnauthorized,
					"message": "InvalidAuthHeader",
					"data":    empty.Empty{},
				},
			)
			c.Abort()
			return
		}

		token := parts[1]
		apiKey, err := m.application.GetApiKeyDetails(c.Request.Context(), token)
		if err != nil || !apiKey.IsActive {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"status":  http.StatusUnauthorized,
					"message": "InvalidAuthHeader",
					"data":    empty.Empty{},
				},
			)
			c.Abort()
			return
		}

		c.Set("application_id", apiKey.ApplicationID)

		c.Next()
	}
}
