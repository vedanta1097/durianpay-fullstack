package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
)

type ListFilter struct {
	Status string
	ID     string
	Sort   string
}

type PaymentRepository interface {
	ListPayments(ctx context.Context, filter ListFilter) ([]entity.Payment, error)
	GetPaymentSummary(ctx context.Context) (entity.PaymentSummary, error)
}

type Payment struct {
	db *sql.DB
}

func NewPaymentRepo(db *sql.DB) *Payment {
	return &Payment{db: db}
}

func (r *Payment) ListPayments(ctx context.Context, filter ListFilter) ([]entity.Payment, error) {
	query := `SELECT id, merchant, status, amount, created_at FROM payments`
	conditions := make([]string, 0, 2)
	args := make([]any, 0, 2)

	if filter.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, filter.Status)
	}
	if filter.ID != "" {
		conditions = append(conditions, "id = ?")
		args = append(args, filter.ID)
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY " + filter.Sort

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, entity.WrapError(err, entity.ErrorCodeInternal, "failed to list payments")
	}
	defer rows.Close()

	payments := make([]entity.Payment, 0)
	for rows.Next() {
		var payment entity.Payment
		if err := rows.Scan(&payment.ID, &payment.Merchant, &payment.Status, &payment.Amount, &payment.CreatedAt); err != nil {
			return nil, entity.WrapError(err, entity.ErrorCodeInternal, "failed to read payment")
		}
		payments = append(payments, payment)
	}
	if err := rows.Err(); err != nil {
		return nil, entity.WrapError(err, entity.ErrorCodeInternal, "failed to list payments")
	}

	return payments, nil
}

func (r *Payment) GetPaymentSummary(ctx context.Context) (entity.PaymentSummary, error) {
	const query = `SELECT
		COUNT(*),
		COALESCE(SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN status = 'processing' THEN 1 ELSE 0 END), 0),
		COALESCE(SUM(CASE WHEN status = 'failed' THEN 1 ELSE 0 END), 0)
	FROM payments`

	var summary entity.PaymentSummary
	err := r.db.QueryRowContext(ctx, query).Scan(
		&summary.Total,
		&summary.Completed,
		&summary.Processing,
		&summary.Failed,
	)
	if err != nil {
		return entity.PaymentSummary{}, entity.WrapError(err, entity.ErrorCodeInternal, "failed to summarize payments")
	}

	return summary, nil
}

func ParseSort(sort string) (string, error) {
	if sort == "" {
		return "created_at DESC, id DESC", nil
	}

	columns := map[string]string{
		"created_at": "created_at",
		"amount":     "amount",
	}
	parts := strings.Split(sort, ",")
	orders := make([]string, 0, len(parts)+1)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		direction := "ASC"
		if strings.HasPrefix(part, "-") {
			direction = "DESC"
			part = strings.TrimPrefix(part, "-")
		}
		column, ok := columns[part]
		if !ok || part == "" {
			return "", errors.New("unsupported sort field")
		}
		orders = append(orders, column+" "+direction)
	}

	orders = append(orders, "id DESC")
	return strings.Join(orders, ", "), nil
}
