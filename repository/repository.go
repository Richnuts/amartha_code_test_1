package repository

import (
	"billing_engine/model"
	"context"
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetLoans(ctx context.Context, exp exp.Expression) ([]model.Loan, error)
	MakePayment(ctx context.Context, input model.Payment) error
	CountPayment(ctx context.Context, loanID string) (int, error)
	UpdateDeliquency() error
}

type repository struct {
	db           *sqlx.DB
	queryBuilder goqu.DialectWrapper
}

func NewRepository(psql *sqlx.DB) Repository {
	return &repository{
		db:           psql,
		queryBuilder: goqu.Dialect("postgres"),
	}
}

func (r *repository) GetLoans(ctx context.Context, filterExp exp.Expression) ([]model.Loan, error) {
	ds := r.queryBuilder.From("loan").Where(filterExp).Order(goqu.I("id").Asc())
	query, params, err := ds.ToSQL()
	if err != nil {
		return nil, err
	}

	var loans []model.Loan
	err = r.db.SelectContext(ctx, &loans, query, params...)
	if err != nil {
		return nil, err
	}

	return loans, nil
}

func (r *repository) MakePayment(ctx context.Context, input model.Payment) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert payment
	ds := r.queryBuilder.From("payment").Insert().Rows(input).Prepared(true)
	query, params, err := ds.ToSQL()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, params...)
	if err != nil {
		return err
	}

	// Update outstanding amount
	preparedStmt := r.queryBuilder.From("loan").Prepared(true)
	updateDs := preparedStmt.
		Update().
		Set(goqu.Record{
			"outstanding_amount": goqu.L("outstanding_amount - ?", input.AmountPaid),
			"last_due_at":        goqu.L("created_at + (interval '1 week' * ?)", input.PaymentNumber+1),
		}).
		Where(goqu.Ex{
			"id": input.LoanID,
		})
	query, params, err = updateDs.ToSQL()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, query, params...)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repository) CountPayment(ctx context.Context, loanID string) (int, error) {
	ds := r.queryBuilder.From("payment").Select(goqu.COUNT(goqu.Star())).Where(goqu.Ex{"loan_id": loanID})
	query, params, err := ds.ToSQL()
	if err != nil {
		return 0, err
	}
	var count int
	err = r.db.GetContext(ctx, &count, query, params...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}

func (r *repository) UpdateDeliquency() error {
	// Update the deliquent status
	preparedStmt := r.queryBuilder.From("loan").Prepared(true)
	updateDs := preparedStmt.
		Update().
		Set(goqu.Record{
			"is_deliquent": true,
		}).
		Where(
			goqu.Ex{
				"is_deliquent": false,
				"outstanding_amount": goqu.Op{
					"neq": 0,
				},
			},
			goqu.L("TO_DATE(TO_TIMESTAMP(NOW()) + INTERVAL '7 hours') > TO_DATE(TO_TIMESTAMP(last_payment_at) + INTERVAL '14 days' + INTERVAL '7 hours')"),
		)
	query, params, err := updateDs.ToSQL()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(query, params...)
	if err != nil {
		return err
	}
	return nil
}
