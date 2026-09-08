package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestListPaymentsFiltersAndSortsAgainstSQLite(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec(`CREATE TABLE payments (id TEXT PRIMARY KEY, merchant TEXT, status TEXT, amount INTEGER, created_at DATETIME)`); err != nil {
		t.Fatal(err)
	}
	for _, payment := range []struct {
		id, status string
		amount     int64
		createdAt  time.Time
	}{
		{"PAY-1", "completed", 200, time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC)},
		{"PAY-2", "completed", 100, time.Date(2026, 3, 3, 0, 0, 0, 0, time.UTC)},
		{"PAY-3", "failed", 300, time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)},
	} {
		if _, err := db.Exec(`INSERT INTO payments(id, merchant, status, amount, created_at) VALUES (?, 'Merchant', ?, ?, ?)`, payment.id, payment.status, payment.amount, payment.createdAt); err != nil {
			t.Fatal(err)
		}
	}

	payments, err := NewPaymentRepo(db).ListPayments(context.Background(), ListFilter{Status: "completed", Sort: "amount ASC, id DESC"})
	if err != nil {
		t.Fatal(err)
	}
	if len(payments) != 2 || payments[0].ID != "PAY-2" || payments[1].ID != "PAY-1" {
		t.Fatalf("payments = %#v", payments)
	}

	summary, err := NewPaymentRepo(db).GetPaymentSummary(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	expected := struct {
		total, completed, processing, failed int
	}{3, 2, 0, 1}
	if summary.Total != expected.total || summary.Completed != expected.completed || summary.Processing != expected.processing || summary.Failed != expected.failed {
		t.Fatalf("summary = %#v", summary)
	}
}

func TestParseSortRejectsUnknownColumn(t *testing.T) {
	if _, err := ParseSort("amount; DROP TABLE payments"); err == nil {
		t.Fatal("ParseSort() accepted an unknown column")
	}
}
