package api

import (
	"fmt"
	"net/http"

	"github.com/asmejia1993/payment-app/async"
	"github.com/asmejia1993/payment-app/bank"
	db "github.com/asmejia1993/payment-app/db/sqlc"
	"github.com/asmejia1993/payment-app/db/util"
	"github.com/asmejia1993/payment-app/token"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Server struct {
	logger     *logrus.Logger
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	plaidMaker bank.Simulator
	router     *gin.Engine
	asynq      async.Queuer
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create token maker: %w", err)
	}

	plaidMaker := bank.NewPlaidClient(config.PlaidClientId, config.PlaidSecret)
	log := NewLogger()
	async := async.NewAsynqClient(config, log)

	server := &Server{
		logger:     log,
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		plaidMaker: plaidMaker,
		asynq:      async,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	router.GET("/monitoring/tasks/*wildcard", server.monitor)
	router.POST("/auth/register", server.createUser)
	router.POST("/auth/login", server.login)
	authRoutes := router.Group("/api/v1").Use(addMiddleware(server.tokenMaker))
	authRoutes.POST("/payments/transaction", server.createTransaction)
	authRoutes.GET("/payments/transaction/:id", server.getTransaction)
	authRoutes.POST("/payments/transaction/:id/refund", server.createTransaction)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
