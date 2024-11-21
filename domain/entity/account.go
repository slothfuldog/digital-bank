package entity

import "time"

type AccountBase struct {
	AccountNumber  string    `json:"account_number"`
	UserID         int       `json:"user_id"`
	AccountType    int       `json:"account_type"`
	BsBalance      float64   `json:"bs_balance"`
	CurrentBalance float64   `json:"current_balance"`
	Status         int       `json:"status"`
	JointWith      int       `json:"joint_with"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type AccountType struct {
	ID          int       `json:"id"`
	AccountName string    `json:"account_name"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
}
