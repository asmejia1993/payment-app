package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	router.GET("/health", server.health)
	router.POST("/auth/register", server.createUser)
	router.POST("/auth/login", server.login)
	authRoutes := router.Group("/api/v1").Use(addMiddleware(server.tokenMaker))
	authRoutes.POST("/payments/transaction", server.createTransaction)
	authRoutes.GET("/payments/transaction/:id", server.getTransaction)
	authRoutes.POST("/payments/transaction/:id/refund", server.createTransaction)

	server.router = router
}

func (server *Server) Run(address string) error {
	srv := &http.Server{
		Addr:    address,
		Handler: server.router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	stopServer := make(chan os.Signal, 1)
	signal.Notify(stopServer, syscall.SIGINT, syscall.SIGTERM)
	<-stopServer
	server.asynq.Close()
	log.Println("server shut down gracefully ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown: %v", err)
	}
	log.Println("Server exiting")
	return nil
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
