package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
)

func (s *Server) monitor(ctx *gin.Context) {
	asynqmon.New(asynqmon.Options{
		RootPath: "/",
		RedisConnOpt: asynq.RedisClientOpt{
			Addr:     s.config.RedisAddr,
			Password: "",
			DB:       0,
		},
	})

}
