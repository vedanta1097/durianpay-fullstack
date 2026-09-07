package usecase

import (
	"context"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
	"github.com/durianpay/fullstack-boilerplate/internal/module/payment/repository"
)

type ListInput struct {
	Status string
	ID     string
	Sort   string
}

type PaymentUsecase interface {
	ListPayments(ctx context.Context, input ListInput) ([]entity.Payment, error)
}

type Payment struct {
	repo repository.PaymentRepository
}

func NewPaymentUsecase(repo repository.PaymentRepository) *Payment {
	return &Payment{repo: repo}
}

func (p *Payment) ListPayments(ctx context.Context, input ListInput) ([]entity.Payment, error) {
	if input.Status != "" && !entity.IsSupportedPaymentStatus(input.Status) {
		return nil, entity.ErrorBadRequest("unsupported payment status")
	}

	sort, err := repository.ParseSort(input.Sort)
	if err != nil {
		return nil, entity.ErrorBadRequest("unsupported sort field")
	}

	return p.repo.ListPayments(ctx, repository.ListFilter{
		Status: input.Status,
		ID:     input.ID,
		Sort:   sort,
	})
}
