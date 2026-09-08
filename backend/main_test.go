package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/durianpay/fullstack-boilerplate/internal/api"
	ah "github.com/durianpay/fullstack-boilerplate/internal/module/auth/handler"
	ar "github.com/durianpay/fullstack-boilerplate/internal/module/auth/repository"
	au "github.com/durianpay/fullstack-boilerplate/internal/module/auth/usecase"
	ph "github.com/durianpay/fullstack-boilerplate/internal/module/payment/handler"
	pr "github.com/durianpay/fullstack-boilerplate/internal/module/payment/repository"
	pu "github.com/durianpay/fullstack-boilerplate/internal/module/payment/usecase"
	"github.com/durianpay/fullstack-boilerplate/internal/openapigen"
	srv "github.com/durianpay/fullstack-boilerplate/internal/service/http"
	_ "github.com/mattn/go-sqlite3"
)

func TestInitDBIsIdempotent(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := initDB(db); err != nil {
		t.Fatal(err)
	}
	if err := initDB(db); err != nil {
		t.Fatal(err)
	}

	for _, table := range []string{"users", "payments"} {
		var count int
		if err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count == 0 {
			t.Fatalf("%s was not seeded", table)
		}
	}
}

func TestLoginThenListPayments(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := initDB(db); err != nil {
		t.Fatal(err)
	}

	authUC := au.NewAuthUsecase(ar.NewUserRepo(db), []byte("test-secret"), time.Hour)
	paymentUC := pu.NewPaymentUsecase(pr.NewPaymentRepo(db))
	server := srv.NewServer(&api.APIHandler{
		Auth:     ah.NewAuthHandler(authUC),
		Payments: ph.NewPaymentHandler(paymentUC),
	}, "../openapi.yaml", authUC)

	loginRequest := httptest.NewRequest(http.MethodPost, "/dashboard/v1/auth/login", bytes.NewBufferString(`{"email":"cs@test.com","password":"password"}`))
	loginRequest.Header.Set("Content-Type", "application/json")
	loginResponse := httptest.NewRecorder()
	server.Routes().ServeHTTP(loginResponse, loginRequest)
	if loginResponse.Code != http.StatusOK {
		t.Fatalf("login status = %d, body = %s", loginResponse.Code, loginResponse.Body.String())
	}

	var login openapigen.LoginResponse
	if err := json.NewDecoder(loginResponse.Body).Decode(&login); err != nil {
		t.Fatal(err)
	}

	paymentRequest := httptest.NewRequest(http.MethodGet, "/dashboard/v1/payments?status=completed", nil)
	paymentRequest.Header.Set("Authorization", "Bearer "+login.Token)
	paymentResponse := httptest.NewRecorder()
	server.Routes().ServeHTTP(paymentResponse, paymentRequest)
	if paymentResponse.Code != http.StatusOK {
		t.Fatalf("payment status = %d, body = %s", paymentResponse.Code, paymentResponse.Body.String())
	}

	var payments openapigen.PaymentListResponse
	if err := json.NewDecoder(paymentResponse.Body).Decode(&payments); err != nil {
		t.Fatal(err)
	}
	if len(payments.Payments) == 0 {
		t.Fatal("completed payment list is empty")
	}
	for _, payment := range payments.Payments {
		if payment.Status != openapigen.PaymentStatusCompleted {
			t.Fatalf("filtered payment status = %q, want %q", payment.Status, openapigen.PaymentStatusCompleted)
		}
	}

	wantSummary := summaryFromPaymentSeeds(paymentSeeds)
	if payments.Summary != wantSummary {
		t.Fatalf("payment summary = %#v, want %#v", payments.Summary, wantSummary)
	}
}

func summaryFromPaymentSeeds(seeds []paymentSeed) openapigen.PaymentSummary {
	summary := openapigen.PaymentSummary{Total: len(seeds)}
	for _, payment := range seeds {
		switch payment.status {
		case "completed":
			summary.Completed++
		case "processing":
			summary.Processing++
		case "failed":
			summary.Failed++
		}
	}
	return summary
}
