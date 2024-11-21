package entity

import "time"

type AccountCodeBase struct {
	AccountCode     string    `json:"account_code"`
	ContraCode      string    `json:"contra_code"`
	CodeType        string    `json:"code_type"`
	CodeDescription string    `json:"code_description"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type LedgerBase struct {
	ID          int       `json:"id"`
	ReferenceNo string    `json:"reference_no"`
	AccountCode string    `json:"account_code"`
	TrxAmount   float64   `json:"trx_amount"`
	GlobalId    string    `json:"global_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LedgerMaster struct {
	ID          int       `json:"id"`
	AccountCode string    `json:"account_code"`
	AccountType string    `json:"account_type"`
	CreditAmt   float64   `json:"credit_amt"`
	DebitAmt    float64   `json:"debit_amt"`
	TotalAmt    float64   `json:"total_amt"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LedgerTransaction struct {
	ID            int       `json:"id"`
	AccountNumber string    `json:"account_number"`
	BeforeBalance float64   `json:"before_balance"`
	TrxAmt        float64   `json:"trx_amt"`
	AfterBalance  float64   `json:"after_balance"`
	GlobalId      string    `json:"global_id"`
	Remark        string    `json:"remark"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `json:"created_by"`
}

type AccrualBase struct {
	ID             int       `json:"id"`
	AccountNumber  string    `json:"account_number"`
	TotalHistory   int       `json:"total_history"`
	BaseAmount     float64   `json:"base_amount"`
	InterestRate   float64   `json:"interest_rate"`
	TotalAccrual   float64   `json:"total_accrual"`
	CurrentAccrued float64   `json:"current_accrued"`
	FromDate       time.Time `json:"from_date"`
	ToDate         time.Time `json:"to_date"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type AccrualHistory struct {
	ID            int       `json:"id"`
	ReferenceNo   string    `json:"reference_no"`
	HistoryNumber int       `json:"history_number"`
	AccrualAmount float64   `json:"accrual_amount"`
	TotalAccrued  float64   `json:"total_accrued"`
	ToBeAccrued   float64   `json:"to_be_accrued"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `json:"created_by"`
}
