package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"go-base-project/internal/middleware"
	"go-base-project/internal/model"
	"go-base-project/internal/service"
)

// ─── Auth Handler ─────────────────────────────────────────────────────────────

type AuthHandler struct {
	jwtSecret   string
	expiryHours int
	log         *logrus.Logger
}

func NewAuthHandler(secret string, expiry int, log *logrus.Logger) *AuthHandler {
	return &AuthHandler{jwtSecret: secret, expiryHours: expiry, log: log}
}

// POST /api/v1/auth/login
// Demo: hardcoded admin/admin — replace with DB lookup in production
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Fail(err.Error()))
		return
	}

	if req.Username != "admin" || req.Password != "admin" {
		c.JSON(http.StatusUnauthorized, model.Fail("invalid credentials"))
		return
	}

	token, err := middleware.GenerateToken(req.Username, h.jwtSecret, h.expiryHours)
	if err != nil {
		h.log.Errorf("generate token: %v", err)
		c.JSON(http.StatusInternalServerError, model.Fail("could not generate token"))
		return
	}

	c.JSON(http.StatusOK, model.OK("login successful", model.LoginResponse{Token: token}))
}

// ─── Customer Handler ─────────────────────────────────────────────────────────

type CustomerHandler struct {
	svc service.CustomerService
	log *logrus.Logger
}

func NewCustomerHandler(svc service.CustomerService, log *logrus.Logger) *CustomerHandler {
	return &CustomerHandler{svc: svc, log: log}
}

// POST /api/v1/customers
func (h *CustomerHandler) Create(c *gin.Context) {
	var req model.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Fail(err.Error()))
		return
	}

	customer, err := h.svc.CreateCustomer(&req)
	if err != nil {
		h.log.Warnf("create customer: %v", err)
		c.JSON(http.StatusConflict, model.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, model.OK("customer created", customer))
}

// GET /api/v1/customers
func (h *CustomerHandler) List(c *gin.Context) {
	customers, err := h.svc.GetAllCustomers()
	if err != nil {
		h.log.Errorf("list customers: %v", err)
		c.JSON(http.StatusInternalServerError, model.Fail("failed to fetch customers"))
		return
	}
	c.JSON(http.StatusOK, model.OK("success", customers))
}

// ─── Transaction Handler ──────────────────────────────────────────────────────

type TransactionHandler struct {
	svc service.TransactionService
	log *logrus.Logger
}

func NewTransactionHandler(svc service.TransactionService, log *logrus.Logger) *TransactionHandler {
	return &TransactionHandler{svc: svc, log: log}
}

// POST /api/v1/transactions
func (h *TransactionHandler) Create(c *gin.Context) {
	var req model.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Fail(err.Error()))
		return
	}

	tx, err := h.svc.CreateTransaction(&req)
	if err != nil {
		h.log.Warnf("create transaction: %v", err)
		c.JSON(http.StatusBadRequest, model.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, model.OK("transaction created", tx))
}

// GET /api/v1/transactions?account_id=1
func (h *TransactionHandler) List(c *gin.Context) {
	type Query struct {
        AccountID int `form:"account_id" binding:"required"`
    }
    var q Query
    if err := c.ShouldBindQuery(&q); err != nil {
        c.JSON(http.StatusBadRequest, model.Fail("account_id is required"))
        return
    }
    txs, err := h.svc.GetByAccount(q.AccountID)
	if err != nil {
		h.log.Errorf("list transactions: %v", err)
		c.JSON(http.StatusInternalServerError, model.Fail("failed to fetch transactions"))
		return
	}
	c.JSON(http.StatusOK, model.OK("success", txs))
}

// GET /api/v1/transactions/passbook?account_id=1&start_date=2017-01-01&end_date=2017-03-15
func (h *TransactionHandler) Passbook(c *gin.Context) {
	var req model.PassbookRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Fail(err.Error()))
		return
	}

	entries, err := h.svc.GetPassbook(&req)
	if err != nil {
		h.log.Warnf("passbook: %v", err)
		c.JSON(http.StatusBadRequest, model.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.OK("success", entries))
}

// ─── Point Handler ────────────────────────────────────────────────────────────

type PointHandler struct {
	svc service.PointService
	log *logrus.Logger
}

func NewPointHandler(svc service.PointService, log *logrus.Logger) *PointHandler {
	return &PointHandler{svc: svc, log: log}
}

// GET /api/v1/points
func (h *PointHandler) List(c *gin.Context) {
	points, err := h.svc.GetCustomerPoints()
	if err != nil {
		h.log.Errorf("get points: %v", err)
		c.JSON(http.StatusInternalServerError, model.Fail("failed to calculate points"))
		return
	}
	c.JSON(http.StatusOK, model.OK("success", points))
}
