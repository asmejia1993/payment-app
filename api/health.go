package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "OK"})
}
