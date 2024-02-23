package db

import "time"

type UserType struct {
	Merchant string
	Customer string
}

var UserTypeEnum UserType = UserType{
	Merchant: "Merchant",
	Customer: "Customer",
}

type User struct {
	Id                string    `json:"id"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	UserType          string    `json:"userType"`
	HashedPassword    string    `json:"hashedPassword"`
	PasswordChangedAt time.Time `json:"passwordChangedAt"`
	CreatedAt         time.Time `json:"createdAt"`
}

type Transaction struct {
	TransactionId string    `json:"transaction_id"`
	RequestId     string    `json:"request_id"`
	Merchant      string    `json:"merchant"`
	Customer      string    `json:"customer"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	Concept       string    `json:"concept"`
	CreatedAt     time.Time `json:"createdAt"`
}
