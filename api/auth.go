package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/asmejia1993/payment-app/async"
	db "github.com/asmejia1993/payment-app/db/sqlc"
	"github.com/asmejia1993/payment-app/db/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	FirstName string `json:"firstName" binding:"required,alphanum"`
	LastName  string `json:"lastName" binding:"required,alphanum"`
	Username  string `json:"username" binding:"required,alphanum"`
	UserType  string `json:"userType" binding:"required,alphanum"`
	Password  string `json:"password" binding:"required,min=6"`
	Email     string `json:"email" binding:"required,email"`
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

type userResponse struct {
	Id                string    `json:"id"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	UserType          string    `json:"userType"`
	PasswordChangedAt time.Time `json:"passwordChangedAt"`
	CreatedAt         time.Time `json:"createdAt"`
}

func NewUserResponse(u db.User) userResponse {
	return userResponse{
		Id:                u.Id,
		FirstName:         u.FirstName,
		LastName:          u.LastName,
		Username:          u.Username,
		UserType:          u.UserType,
		Email:             u.Email,
		PasswordChangedAt: u.PasswordChangedAt,
		CreatedAt:         u.CreatedAt,
	}
}

func (s *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		s.logger.Errorf("error binding the request: %v", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Username:       req.Username,
		UserType:       req.UserType,
		HashedPassword: hashedPassword,
		Email:          req.Email,
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				s.logger.Errorf("unique violation error: %v", err)
				ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
				return
			}
		}
		s.logger.Errorf("error creating an user: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	eventLog := async.AuditLogEntry{
		Actor:   user.Email,
		Action:  "POST",
		Module:  "Auth",
		When:    time.Now(),
		Details: "new user on boarding",
	}
	err = s.asynq.Enqueue(eventLog, async.TypeNewUser)
	if err != nil {
		s.logger.Errorf("register: %v", err)
	}

	resp := NewUserResponse(user)
	ctx.JSON(http.StatusCreated, resp)
}

func (s *Server) login(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		s.logger.Errorf("error binding the request: %v", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.store.GetUser(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err := util.CheckPassword(user.HashedPassword, req.Password); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, _, err := s.tokenMaker.CreateToken(user.Email, s.config.AccessTokenDuration)
	if err != nil {
		s.logger.Errorf("error creating a token: %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	eventLog := async.AuditLogEntry{
		Actor:   user.Email,
		Action:  "POST",
		Module:  "Auth",
		When:    time.Now(),
		Details: "login",
	}
	err = s.asynq.Enqueue(eventLog, async.TypeLoginUser)
	if err != nil {
		s.logger.Errorf("login: %v", err)
	}

	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        NewUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}
