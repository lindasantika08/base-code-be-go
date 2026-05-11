package service

import (
	"fmt"
	"math"
	"time"

	"go-base-project/internal/model"
	"go-base-project/internal/repository"
)

// ─── Customer Service ─────────────────────────────────────────────────────────

type CustomerService interface {
	CreateCustomer(req *model.CreateCustomerRequest) (*model.Customer, error)
	GetAllCustomers() ([]model.Customer, error)
}

type customerSvc struct{ repo repository.CustomerRepository }

func NewCustomerService(r repository.CustomerRepository) CustomerService {
	return &customerSvc{repo: r}
}

func (s *customerSvc) CreateCustomer(req *model.CreateCustomerRequest) (*model.Customer, error) {
	// Check duplicate
	existing, err := s.repo.FindByID(req.AccountID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("account_id %d already exists", req.AccountID)
	}

	c := &model.Customer{AccountID: req.AccountID, Name: req.Name}
	if err := s.repo.Create(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *customerSvc) GetAllCustomers() ([]model.Customer, error) {
	return s.repo.FindAll()
}

// ─── Transaction Service ──────────────────────────────────────────────────────

type TransactionService interface {
	CreateTransaction(req *model.CreateTransactionRequest) (*model.Transaction, error)
	GetByAccount(accountID int) ([]model.Transaction, error)
	GetPassbook(req *model.PassbookRequest) ([]model.PassbookEntry, error)
}

type transactionSvc struct {
	txRepo   repository.TransactionRepository
	custRepo repository.CustomerRepository
}

func NewTransactionService(tr repository.TransactionRepository, cr repository.CustomerRepository) TransactionService {
	return &transactionSvc{txRepo: tr, custRepo: cr}
}

func (s *transactionSvc) CreateTransaction(req *model.CreateTransactionRequest) (*model.Transaction, error) {
	// Validate account exists
	cust, err := s.custRepo.FindByID(req.AccountID)
	if err != nil {
		return nil, err
	}
	if cust == nil {
		return nil, fmt.Errorf("account_id %d not found", req.AccountID)
	}

	date, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction_date format, use YYYY-MM-DD")
	}

	t := &model.Transaction{
		AccountID:         req.AccountID,
		TransactionDate:   date,
		Description:       req.Description,
		DebitCreditStatus: req.DebitCreditStatus,
		Amount:            req.Amount,
	}
	if err := s.txRepo.Create(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *transactionSvc) GetByAccount(accountID int) ([]model.Transaction, error) {
	return s.txRepo.FindByAccount(accountID)
}

func (s *transactionSvc) GetPassbook(req *model.PassbookRequest) ([]model.PassbookEntry, error) {
	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format")
	}
	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date format")
	}
	// Include full end day
	end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	txs, err := s.txRepo.FindByAccountAndDateRange(req.AccountID, start, end)
	if err != nil {
		return nil, err
	}

	var entries []model.PassbookEntry
	for _, t := range txs {
		e := model.PassbookEntry{
			TransactionDate: t.TransactionDate,
			Description:     t.Description,
			Amount:          t.Amount,
		}
		if t.DebitCreditStatus == "C" {
			e.Credit = &t.Amount
		} else {
			e.Debit = &t.Amount
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// ─── Point Service ────────────────────────────────────────────────────────────

type PointService interface {
	GetCustomerPoints() ([]model.CustomerPoint, error)
}

type pointSvc struct {
	custRepo repository.CustomerRepository
	txRepo   repository.TransactionRepository
}

func NewPointService(cr repository.CustomerRepository, tr repository.TransactionRepository) PointService {
	return &pointSvc{custRepo: cr, txRepo: tr}
}

func (s *pointSvc) GetCustomerPoints() ([]model.CustomerPoint, error) {
	customers, err := s.custRepo.FindAll()
	if err != nil {
		return nil, err
	}

	txs, err := s.txRepo.FindAllForPoints()
	if err != nil {
		return nil, err
	}

	// Map accountID -> total points
	pointMap := make(map[int]float64)
	for _, t := range txs {
		var p float64
		switch t.Description {
		case "Beli Pulsa":
			p = calcPulsaPoint(t.Amount)
		case "Bayar Listrik":
			p = calcListrikPoint(t.Amount)
		}
		pointMap[t.AccountID] += p
	}

	var result []model.CustomerPoint
	for _, c := range customers {
		result = append(result, model.CustomerPoint{
			AccountID:  c.AccountID,
			Name:       c.Name,
			TotalPoint: pointMap[c.AccountID],
		})
	}
	return result, nil
}

// calcPulsaPoint: 0-10000 → 0pt | 10001-30000 → 1pt/1000 | >30000 → 2pt/1000
func calcPulsaPoint(amount float64) float64 {
	var total float64
	if amount <= 10000 {
		return 0
	}
	if amount > 10000 {
		tier2 := math.Min(amount, 30000) - 10000
		total += math.Floor(tier2/1000) * 1
	}
	if amount > 30000 {
		tier3 := amount - 30000
		total += math.Floor(tier3/1000) * 2
	}
	return total
}

// calcListrikPoint: 0-50000 → 0pt | 50001-100000 → 1pt/2000 | >100000 → 2pt/2000
func calcListrikPoint(amount float64) float64 {
	var total float64
	if amount <= 50000 {
		return 0
	}
	if amount > 50000 {
		tier2 := math.Min(amount, 100000) - 50000
		total += math.Floor(tier2/2000) * 1
	}
	if amount > 100000 {
		tier3 := amount - 100000
		total += math.Floor(tier3/2000) * 2
	}
	return total
}
