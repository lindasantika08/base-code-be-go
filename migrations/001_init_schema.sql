CREATE DATABASE IF NOT EXISTS banking_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE banking_db;

-- ─── Customers ────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS customers (
    account_id  INT         NOT NULL,
    name        VARCHAR(100) NOT NULL,
    created_at  DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT pk_customers PRIMARY KEY (account_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ─── Transactions ─────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS transactions (
    id                   INT           NOT NULL AUTO_INCREMENT,
    account_id           INT           NOT NULL,
    transaction_date     DATETIME      NOT NULL,
    description          VARCHAR(255)  NOT NULL,
    debit_credit_status  CHAR(1)       NOT NULL COMMENT 'D=Debit, C=Credit',
    amount               DECIMAL(15,2) NOT NULL,
    created_at           DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pk_transactions     PRIMARY KEY (id),
    CONSTRAINT fk_transactions_customer FOREIGN KEY (account_id)
        REFERENCES customers (account_id) ON DELETE RESTRICT,
    CONSTRAINT chk_dc_status CHECK (debit_credit_status IN ('D','C')),
    CONSTRAINT chk_amount    CHECK (amount > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Index for common queries
CREATE INDEX idx_transactions_account    ON transactions (account_id);
CREATE INDEX idx_transactions_date       ON transactions (transaction_date);
CREATE INDEX idx_transactions_account_dt ON transactions (account_id, transaction_date);
CREATE INDEX idx_transactions_desc       ON transactions (description);

-- ─── Seed data (from soal) ────────────────────────────────────────────────────
INSERT INTO customers (account_id, name) VALUES
  (1, 'Customer1'),
  (2, 'Customer2');

INSERT INTO transactions (account_id, transaction_date, description, debit_credit_status, amount) VALUES
  (1, '2017-01-01', 'Setor Tunai',   'C', 200000),
  (1, '2017-01-05', 'Beli Pulsa',    'D',  10000),
  (1, '2017-01-06', 'Bayar Listrik', 'D',  70000),
  (1, '2017-01-07', 'Tarik Tunai',   'D', 100000),
  (1, '2017-02-01', 'Setor Tunai',   'C', 300000),
  (1, '2017-02-05', 'Bayar Listrik', 'D',  50000),
  (1, '2017-02-15', 'Tarik Tunai',   'D',  50000),
  (1, '2017-02-20', 'Beli Pulsa',    'D',  40000),
  (1, '2017-02-28', 'Tarik Tunai',   'D',  50000),
  (1, '2017-03-01', 'Setor Tunai',   'C',  50000),
  (1, '2017-03-07', 'Bayar Listrik', 'D', 125000),
  (1, '2017-03-15', 'Beli Pulsa',    'D',  20000);
