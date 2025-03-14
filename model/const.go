package model

type DurationUnit string

const (
	WEEK DurationUnit = "WEEK"
)

// error const
const (
	ISE_MESSAGE     = "Internal server error"
	INVALID_LOAN    = "invalid loan"
	LOAN_PAID       = "loan is paid"
	INVALID_PAYMENT = "invalid payment please pay %v instead of %v"
)
