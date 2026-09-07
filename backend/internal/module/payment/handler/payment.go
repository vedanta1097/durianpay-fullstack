package handler

import (
	"encoding/json"
	"net/http"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
	"github.com/durianpay/fullstack-boilerplate/internal/module/payment/usecase"
	"github.com/durianpay/fullstack-boilerplate/internal/openapigen"
	"github.com/durianpay/fullstack-boilerplate/internal/transport"
)

type PaymentHandler struct {
	paymentUC usecase.PaymentUsecase
}

func NewPaymentHandler(paymentUC usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{paymentUC: paymentUC}
}

func (h *PaymentHandler) GetDashboardV1Payments(w http.ResponseWriter, r *http.Request, params openapigen.GetDashboardV1PaymentsParams) {
	input := usecase.ListInput{}
	if params.Status != nil {
		input.Status = string(*params.Status)
	}
	if params.Id != nil {
		input.ID = *params.Id
	}
	if params.Sort != nil {
		input.Sort = string(*params.Sort)
	}

	payments, err := h.paymentUC.ListPayments(r.Context(), input)
	if err != nil {
		transport.WriteError(w, err)
		return
	}

	responsePayments := make([]openapigen.Payment, 0, len(payments))
	for _, payment := range payments {
		responsePayments = append(responsePayments, openapigen.Payment{
			Id:        payment.ID,
			Merchant:  payment.Merchant,
			Status:    openapigen.PaymentStatus(payment.Status),
			Amount:    payment.Amount,
			CreatedAt: payment.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(openapigen.PaymentListResponse{Payments: &responsePayments}); err != nil {
		transport.WriteAppError(w, entity.ErrorInternal("internal server error"))
	}
}
