package model

import (
	"time"

	"github.com/google/uuid"
)

type Loan struct {
	ID                uuid.UUID    `db:"id" json:"id"`
	UserID            uuid.UUID    `db:"user_id" json:"userId"`
	PrincipalAmount   int          `db:"principal_amount" json:"principalAmount"`
	InterestRate      float64      `db:"interest_rate" json:"interestRate"`
	OutstandingAmount int          `db:"outstanding_amount" json:"outstandingAmount"`
	Duration          int          `db:"duration" json:"duration"`
	DurationUnit      DurationUnit `db:"duration_unit" json:"durationUnit"`
	CreatedAt         time.Time    `db:"created_at" json:"createdAt"`
	IsDeliquent       bool         `db:"is_deliquent" json:"isDeliquent"`
	LastDueAt         time.Time    `db:"last_due_at" json:"lastDueAt"`
}

type Payment struct {
	ID            uuid.UUID `db:"id" json:"id"`
	LoanID        string    `db:"loan_id" json:"loanId"`
	PaymentNumber int       `db:"payment_number" json:"paymentNumber"`
	AmountPaid    int       `db:"amount_paid" json:"amountPaid" validate:"required"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
}

type LoanSchedule struct {
	PrincipalAmount int                  `json:"principalAmount"`
	InterestRate    float64              `json:"interestRate"`
	TotalAmount     int                  `json:"TotalAmount"`
	Duration        int                  `json:"duration"`
	DurationUnit    DurationUnit         `json:"durationUnit"`
	Details         []LoanScheduleDetail `json:"details"`
}

type LoanScheduleDetail struct {
	PaymentAmount int `json:"paymentAmount"`
	PaymentNumber int `json:"paymentNumber"`
}

func (ths *Payment) ConstructPayment(paymentNumber int) {
	ths.CreatedAt = time.Now()
	ths.ID = uuid.New()
	ths.PaymentNumber = paymentNumber
}
