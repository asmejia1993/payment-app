package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/asmejia1993/payment-app/async"
	db "github.com/asmejia1993/payment-app/db/sqlc"
	"github.com/asmejia1993/payment-app/token"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type paymentRequest struct {
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	Email        string  `json:"email"`
	ToMerchant   string  `json:"merchant"`
	Card         string  `json:"card"`
	Concept      string  `json:"concept"`
	FromCustomer string  `json:"customer"`
}

type successResponse struct {
	Success   string `json:"success"`
	RequestId string `json:"request_id"`
}

func (s *Server) createTransaction(ctx *gin.Context) {
	var req paymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		s.logger.Errorf("error binding the request: %v", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Email != req.Email {
		err := errors.New("user is not authorized to make the process")
		s.logger.Errorf("unauthorized: %v", err)
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	transactionID := uuid.New().String()
	s.logger.Infof("transaction id: %s", transactionID)

	result, err := s.plaidMaker.StartPaymentProcess(ctx)
	if err != nil {
		s.logger.Errorf("error in payment process: %v", err)
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	s.logger.Infof("result from payment process: %s\n", result)

	trx := db.CreateTransactionHistory{
		TransactionId: transactionID,
		RequestId:     uuid.New().String(),
		Merchant:      req.ToMerchant,
		Customer:      req.FromCustomer,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Concept:       req.Concept,
	}

	s.logger.Infof("transaction to insert: %v", trx)

	res, err := s.store.CreateTransactionHistory(ctx, trx)
	if err != nil {
		s.logger.Errorf("error creating a transaction: %v", err)
	}

	eventLog := async.AuditLogEntry{
		Actor:   trx.TransactionId,
		Action:  "POST",
		Module:  "Payment",
		When:    time.Now(),
		Details: "creating a new transaction",
	}
	err = s.asynq.Enqueue(eventLog, async.TypeNewTransaction)
	if err != nil {
		s.logger.Errorf("creating a transaction: %v", err)
	}

	ctx.JSON(http.StatusCreated, successResponse{Success: "ok", RequestId: res.TransactionId})
}

func (s *Server) getTransaction(ctx *gin.Context) {
	transactionId, _ := ctx.Params.Get("id")
	token := strings.Split(ctx.Request.Header["Authorization"][0], " ")[1]

	if len(token) <= 0 {
		err := errors.New("user is not authorized to make the process")
		s.logger.Errorf("unauthorized: %v", err)
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, err := s.tokenMaker.VerifyToken(token)
	if err != nil {
		err := errors.New("token invalid")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	eventLog := async.AuditLogEntry{
		Actor:   transactionId,
		Action:  "GET",
		Module:  "Payment",
		When:    time.Now(),
		Details: "getting a transaction detail",
	}
	err = s.asynq.Enqueue(eventLog, async.TypeNewTransaction)
	if err != nil {
		s.logger.Errorf("sending a transaction: %v", err)
	}

	res, _ := s.store.GetTransactionDetails(ctx, transactionId)
	ctx.JSON(http.StatusOK, gin.H{"data": res})
}
