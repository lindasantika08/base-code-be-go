package model

import "time"

// ─── Nasabah ──────────────────────────────────────────────────────────────────

type Customer struct {
	AccountID int    `json:"account_id" db:"account_id"`
	Name      string `json:"name"       db:"name"`
}

type CreateCustomerRequest struct {
	AccountID int    `json:"account_id" binding:"required"`
	Name      string `json:"name"       binding:"required,min=1,max=100"`
}

// ─── Transaksi ────────────────────────────────────────────────────────────────

type Transaction struct {
	ID                int       `json:"id"                  db:"id"`
	AccountID         int       `json:"account_id"          db:"account_id"`
	TransactionDate   time.Time `json:"transaction_date"    db:"transaction_date"`
	Description       string    `json:"description"         db:"description"`
	DebitCreditStatus string    `json:"debit_credit_status" db:"debit_credit_status"` // 'D' | 'C'
	Amount            float64   `json:"amount"              db:"amount"`
}

type CreateTransactionRequest struct {
	AccountID         int     `json:"account_id"          binding:"required"`
	TransactionDate   string  `json:"transaction_date"    binding:"required"` // "2006-01-02"
	Description       string  `json:"description"         binding:"required,min=1,max=255"`
	DebitCreditStatus string  `json:"debit_credit_status" binding:"required,oneof=D C"`
	Amount            float64 `json:"amount"              binding:"required,gt=0"`
}

// ─── Point ────────────────────────────────────────────────────────────────────

type CustomerPoint struct {
	AccountID  int     `json:"account_id"   db:"account_id"`
	Name       string  `json:"name"         db:"name"`
	TotalPoint float64 `json:"total_point"  db:"total_point"`
}

// ─── Buku Tabungan ────────────────────────────────────────────────────────────

type PassbookRequest struct {
	AccountID int    `form:"account_id" binding:"required"`
	StartDate string `form:"start_date" binding:"required"` // "2006-01-02"
	EndDate   string `form:"end_date"   binding:"required"` // "2006-01-02"
}

type PassbookEntry struct {
	TransactionDate   time.Time `json:"transaction_date"    db:"transaction_date"`
	Description       string    `json:"description"         db:"description"`
	Credit            *float64  `json:"credit"              db:"credit"`
	Debit             *float64  `json:"debit"               db:"debit"`
	Amount            float64   `json:"amount"              db:"amount"`
}

// ─── Auth ─────────────────────────────────────────────────────────────────────

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// ─── Generic Response ─────────────────────────────────────────────────────────

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func OK(msg string, data interface{}) Response {
	return Response{Success: true, Message: msg, Data: data}
}

func Fail(msg string) Response {
	return Response{Success: false, Message: msg}
}
