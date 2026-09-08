package usecase

import (
	"context"
	"testing"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
	"github.com/durianpay/fullstack-boilerplate/internal/module/payment/repository"
)

type paymentRepositoryStub struct {
	filter  repository.ListFilter
	called  bool
	summary entity.PaymentSummary
}

func (s *paymentRepositoryStub) ListPayments(_ context.Context, filter repository.ListFilter) ([]entity.Payment, error) {
	s.called = true
	s.filter = filter
	return []entity.Payment{}, nil
}

func (s *paymentRepositoryStub) GetPaymentSummary(_ context.Context) (entity.PaymentSummary, error) {
	return s.summary, nil
}

func TestListPaymentsPassesValidatedFilterToRepository(t *testing.T) {
	repo := &paymentRepositoryStub{summary: entity.PaymentSummary{Total: 30, Completed: 13, Processing: 9, Failed: 8}}
	usecase := NewPaymentUsecase(repo)

	result, err := usecase.ListPayments(context.Background(), ListInput{Status: entity.PaymentStatusCompleted, ID: "PAY-0001", Sort: "-amount,created_at"})
	if err != nil {
		t.Fatalf("ListPayments() error = %v", err)
	}
	if !repo.called {
		t.Fatal("repository was not called")
	}
	if repo.filter.Status != entity.PaymentStatusCompleted || repo.filter.ID != "PAY-0001" {
		t.Fatalf("filter = %#v", repo.filter)
	}
	if repo.filter.Sort != "amount DESC, created_at ASC, id DESC" {
		t.Fatalf("sort = %q", repo.filter.Sort)
	}
	if result.Summary != repo.summary {
		t.Fatalf("summary = %#v", result.Summary)
	}
}

func TestListPaymentsRejectsUnsupportedInputBeforeRepository(t *testing.T) {
	for _, input := range []ListInput{{Status: "cancelled"}, {Sort: "merchant"}} {
		repo := &paymentRepositoryStub{}
		_, err := NewPaymentUsecase(repo).ListPayments(context.Background(), input)
		if err == nil {
			t.Fatal("ListPayments() error = nil")
		}
		if repo.called {
			t.Fatal("repository was called for invalid input")
		}
	}
}
