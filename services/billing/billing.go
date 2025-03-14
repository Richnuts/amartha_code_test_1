package billing

import (
	"billing_engine/model"
	"billing_engine/repository"
	"billing_engine/utils"
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Service interface {
	GetOutstanding(c echo.Context) error
	IsDelinquent(c echo.Context) error
	MakePayment(c echo.Context) error
	GetLoanSchedule(c echo.Context) error
	//daily job
	UpdateDeliquency() error
}

type service struct {
	repo repository.Repository
}

func NewService(
	repo repository.Repository,
) Service {
	return &service{
		repo: repo,
	}
}

func (ths *service) GetOutstanding(c echo.Context) error {
	ctx := c.Request().Context()
	loanID := c.Param("loan_id")
	_, err := uuid.Parse(loanID)
	if err != nil {
		return utils.InvalidRequest(err)
	}
	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.InvalidRequest(err)
	}
	res, err := ths.repo.GetLoans(ctx, goqu.Ex{
		"user_id": userID,
		"id":      loanID,
	})
	if err != nil {
		return utils.InternalServerError(utils.GetRequestID(c), err)
	} else if len(res) == 0 {
		return utils.InvalidRequest(errors.New(model.INVALID_LOAN))
	}
	return c.JSON(http.StatusOK, utils.Response(res[0]))
}

func (ths *service) IsDelinquent(c echo.Context) error {
	ctx := c.Request().Context()
	loanID := c.Param("loan_id")
	_, err := uuid.Parse(loanID)
	if err != nil {
		return utils.InvalidRequest(err)
	}
	// get the payment and check if the returned duration match the needed
	loans, err := ths.repo.GetLoans(ctx, goqu.Ex{
		"id": loanID,
	})
	if err != nil {
		return utils.InternalServerError(utils.GetRequestID(c), err)
	} else if len(loans) == 0 {
		return utils.InvalidRequest(errors.New(model.INVALID_LOAN))
	}
	loan := loans[0]

	// durasi pembayaran harusnya di berapa ? check dari due terakhir
	timeConverted := loan.LastDueAt.Local().Truncate(time.Hour)

	// bandingin kalo hasilnya total pembayaran + 2 kurang dari seharusnya jadi deliquent
	return c.JSON(http.StatusOK, map[string]bool{
		"isDeliquent": time.Now().After(timeConverted.Add(time.Hour * 24 * 7)),
	})
}

func (ths *service) MakePayment(c echo.Context) error {
	ctx := c.Request().Context()
	input, err := utils.BindAndValidateGeneric[model.Payment](c)
	if err != nil {
		return utils.InvalidRequest(err)
	}
	userID, err := utils.GetUserID(c)
	if err != nil {
		return utils.InvalidRequest(err)
	}
	input.LoanID = c.Param("loan_id")
	_, err = uuid.Parse(input.LoanID)
	if err != nil {
		return utils.InvalidRequest(err)
	}
	// validate the loan
	loans, err := ths.repo.GetLoans(ctx, goqu.Ex{
		"id":      input.LoanID,
		"user_id": userID,
	})
	if err != nil {
		return utils.InternalServerError(utils.GetRequestID(c), err)
	} else if len(loans) == 0 {
		return utils.InvalidRequest(errors.New(model.INVALID_LOAN))
	}

	if loans[0].OutstandingAmount == 0 {
		return utils.InvalidRequest(errors.New(model.LOAN_PAID))
	}

	// paymentNumber, err := ths.repo.CountPayment(ctx, input.LoanID)
	// if err != nil {
	// 	return utils.InternalServerError(utils.GetRequestID(c), err)
	// }
	paymentNumber := 51
	input.ConstructPayment(paymentNumber + 1)

	if paymentNumber+1 == loans[0].Duration {
		if (loans[0].OutstandingAmount - input.AmountPaid) != 0 {
			return utils.InvalidRequest(fmt.Errorf(model.INVALID_PAYMENT, loans[0].OutstandingAmount, input.AmountPaid))
		}
	}

	err = ths.repo.MakePayment(ctx, input)
	if err != nil {
		return utils.InternalServerError(utils.GetRequestID(c), err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "success",
	})
}

func (ths *service) GetLoanSchedule(c echo.Context) error {
	ctx := c.Request().Context()
	loanID := c.Param("loan_id")
	_, err := uuid.Parse(loanID)
	if err != nil {
		return utils.InvalidRequest(err)
	}

	loans, err := ths.repo.GetLoans(ctx, goqu.Ex{
		"id": loanID,
	})
	if err != nil {
		return utils.InternalServerError(utils.GetRequestID(c), err)
	}

	loan := loans[0]
	res := model.LoanSchedule{
		Duration:        loan.Duration,
		PrincipalAmount: loan.PrincipalAmount,
		TotalAmount:     int(math.Round(float64(loan.PrincipalAmount) * (1 + loan.InterestRate/100.0))),
		InterestRate:    loan.InterestRate,
		DurationUnit:    loan.DurationUnit,
		Details:         make([]model.LoanScheduleDetail, loan.Duration),
	}

	paymentAmount := res.TotalAmount / loan.Duration
	totalPayment := 0
	for i := range loan.Duration {
		if i == loan.Duration-1 {
			paymentAmount = res.TotalAmount - totalPayment
		}

		res.Details[i] = model.LoanScheduleDetail{
			PaymentAmount: paymentAmount,
			PaymentNumber: i + 1,
		}

		totalPayment += paymentAmount
	}

	return c.JSON(http.StatusOK, utils.Response(res))
}

func (ths *service) UpdateDeliquency() error {
	return ths.repo.UpdateDeliquency()
}
