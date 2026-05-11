package repository

import (
	"database/sql"
	"fmt"
	"time"

	"go-base-project/internal/model"
)

type CustomerRepository interface {
	Create(c *model.Customer) error
	FindAll() ([]model.Customer, error)
	FindByID(id int) (*model.Customer, error)
}

type TransactionRepository interface {
	Create(t *model.Transaction) error
	FindByAccount(accountID int) ([]model.Transaction, error)
	FindByAccountAndDateRange(accountID int, start, end time.Time) ([]model.Transaction, error)
	FindAllForPoints() ([]model.Transaction, error)
}

// ─── Customer ─────────────────────────────────────────────────────────────────

type customerRepo struct{ db *sql.DB }

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &customerRepo{db: db}
}

func (r *customerRepo) Create(c *model.Customer) error {
	_, err := r.db.Exec(
		`INSERT INTO customers (account_id, name) VALUES (?, ?)`,
		c.AccountID, c.Name,
	)
	if err != nil {
		return fmt.Errorf("create customer: %w", err)
	}
	return nil
}

func (r *customerRepo) FindAll() ([]model.Customer, error) {
	rows, err := r.db.Query(`SELECT account_id, name FROM customers ORDER BY account_id`)
	if err != nil {
		return nil, fmt.Errorf("find all customers: %w", err)
	}
	defer rows.Close()

	var list []model.Customer
	for rows.Next() {
		var c model.Customer
		if err := rows.Scan(&c.AccountID, &c.Name); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

func (r *customerRepo) FindByID(id int) (*model.Customer, error) {
	var c model.Customer
	err := r.db.QueryRow(
		`SELECT account_id, name FROM customers WHERE account_id = ?`, id,
	).Scan(&c.AccountID, &c.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &c, err
}

// ─── Transaction ──────────────────────────────────────────────────────────────

type transactionRepo struct{ db *sql.DB }

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) Create(t *model.Transaction) error {
	res, err := r.db.Exec(
		`INSERT INTO transactions (account_id, transaction_date, description, debit_credit_status, amount)
		 VALUES (?, ?, ?, ?, ?)`,
		t.AccountID, t.TransactionDate, t.Description, t.DebitCreditStatus, t.Amount,
	)
	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}
	id, _ := res.LastInsertId()
	t.ID = int(id)
	return nil
}

func (r *transactionRepo) FindByAccount(accountID int) ([]model.Transaction, error) {
	rows, err := r.db.Query(
		`SELECT id, account_id, transaction_date, description, debit_credit_status, amount
		 FROM transactions WHERE account_id = ? ORDER BY transaction_date`,
		accountID,
	)
	if err != nil {
		return nil, fmt.Errorf("find transactions: %w", err)
	}
	defer rows.Close()
	return scanTransactions(rows)
}

func (r *transactionRepo) FindByAccountAndDateRange(accountID int, start, end time.Time) ([]model.Transaction, error) {
	rows, err := r.db.Query(
		`SELECT id, account_id, transaction_date, description, debit_credit_status, amount
		 FROM transactions
		 WHERE account_id = ? AND transaction_date BETWEEN ? AND ?
		 ORDER BY transaction_date`,
		accountID, start, end,
	)
	if err != nil {
		return nil, fmt.Errorf("find transactions by range: %w", err)
	}
	defer rows.Close()
	return scanTransactions(rows)
}

func (r *transactionRepo) FindAllForPoints() ([]model.Transaction, error) {
	rows, err := r.db.Query(
		`SELECT id, account_id, transaction_date, description, debit_credit_status, amount
		 FROM transactions
		 WHERE description IN ('Beli Pulsa','Bayar Listrik')`,
	)
	if err != nil {
		return nil, fmt.Errorf("find point transactions: %w", err)
	}
	defer rows.Close()
	return scanTransactions(rows)
}

func scanTransactions(rows *sql.Rows) ([]model.Transaction, error) {
	var list []model.Transaction
	for rows.Next() {
		var t model.Transaction
		if err := rows.Scan(&t.ID, &t.AccountID, &t.TransactionDate,
			&t.Description, &t.DebitCreditStatus, &t.Amount); err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, rows.Err()
}
