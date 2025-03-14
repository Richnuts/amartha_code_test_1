package route

import (
	"billing_engine/services/billing"

	"github.com/labstack/echo/v4"
)

type Services struct {
	BillingService billing.Service
}

func NewRoute(e *echo.Echo, services Services) {
	e.GET("/billing/outstanding/:loan_id", services.BillingService.GetOutstanding)
	e.GET("/billing/deliquent/:loan_id", services.BillingService.IsDelinquent)
	e.GET("/billing/loan-schedule/:loan_id", services.BillingService.GetLoanSchedule)
	e.POST("/billing/:loan_id", services.BillingService.MakePayment)
}
