package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *RouteHandler) crossOriginResourceSharingHandler(c *gin.Context) {
	// Always return 200
	c.Header("Access-Control-Allow-Origin", c.Request.Host)
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.AbortWithStatusJSON(http.StatusOK, returnJSON)
}
