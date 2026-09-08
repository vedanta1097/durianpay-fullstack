package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/durianpay/fullstack-boilerplate/internal/api"
	"github.com/durianpay/fullstack-boilerplate/internal/config"
	ah "github.com/durianpay/fullstack-boilerplate/internal/module/auth/handler"
	ar "github.com/durianpay/fullstack-boilerplate/internal/module/auth/repository"
	au "github.com/durianpay/fullstack-boilerplate/internal/module/auth/usecase"
	ph "github.com/durianpay/fullstack-boilerplate/internal/module/payment/handler"
	pr "github.com/durianpay/fullstack-boilerplate/internal/module/payment/repository"
	pu "github.com/durianpay/fullstack-boilerplate/internal/module/payment/usecase"
	srv "github.com/durianpay/fullstack-boilerplate/internal/service/http"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	_ = godotenv.Load()

	db, err := sql.Open("sqlite3", config.DatabasePath+"?_foreign_keys=1")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := initDB(db); err != nil {
		log.Fatal(err)
	}

	JwtExpiredDuration, err := time.ParseDuration(config.JwtExpired)
	if err != nil {
		panic(err)
	}

	userRepo := ar.NewUserRepo(db)

	authUC := au.NewAuthUsecase(userRepo, config.JwtSecret, JwtExpiredDuration)

	authH := ah.NewAuthHandler(authUC)
	paymentRepo := pr.NewPaymentRepo(db)
	paymentUC := pu.NewPaymentUsecase(paymentRepo)
	paymentH := ph.NewPaymentHandler(paymentUC)

	apiHandler := &api.APIHandler{
		Auth:     authH,
		Payments: paymentH,
	}

	server := srv.NewServer(apiHandler, config.OpenapiYamlLocation, authUC)

	addr := config.HttpAddress
	log.Printf("starting server on %s", addr)
	server.Start(addr)
}

func initDB(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
		  id INTEGER PRIMARY KEY AUTOINCREMENT,
		  email TEXT NOT NULL UNIQUE,
		  password_hash TEXT NOT NULL,
		  role TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS payments (
		  id TEXT PRIMARY KEY,
		  merchant TEXT NOT NULL,
		  status TEXT NOT NULL CHECK (status IN ('completed', 'processing', 'failed')),
		  amount INTEGER NOT NULL CHECK (amount >= 0),
		  created_at DATETIME NOT NULL
		);`,
		`CREATE INDEX IF NOT EXISTS idx_payments_created_at_id ON payments(created_at DESC, id DESC);`,
		`CREATE INDEX IF NOT EXISTS idx_payments_status_created_at_id ON payments(status, created_at DESC, id DESC);`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return err
		}
	}
	if err := seedDB(db); err != nil {
		return err
	}

	const dbLifetime = time.Minute * 5
	db.SetConnMaxLifetime(dbLifetime)
	return nil
}

func seedDB(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	for _, user := range []struct{ email, role string }{
		{"cs@test.com", "cs"},
		{"operation@test.com", "operation"},
	} {
		if _, err := tx.Exec(`INSERT INTO users(email, password_hash, role) VALUES (?, ?, ?) ON CONFLICT(email) DO NOTHING`, user.email, string(hash), user.role); err != nil {
			return err
		}
	}

	for _, payment := range paymentSeeds {
		if _, err := tx.Exec(`INSERT INTO payments(id, merchant, status, amount, created_at) VALUES (?, ?, ?, ?, ?) ON CONFLICT(id) DO NOTHING`, payment.id, payment.merchant, payment.status, payment.amount, payment.createdAt); err != nil {
			return err
		}
	}

	return tx.Commit()
}

type paymentSeed struct {
	id        string
	merchant  string
	status    string
	amount    int64
	createdAt string
}

var paymentSeeds = []paymentSeed{
	{"PAY-0001", "Kopi Nusantara", "completed", 125000, "2026-03-10T08:15:00Z"},
	{"PAY-0002", "Toko Sejahtera", "processing", 89000, "2026-03-10T07:30:00Z"},
	{"PAY-0003", "Batik Indah", "failed", 210000, "2026-03-09T15:45:00Z"},
	{"PAY-0004", "Warung Rasa", "completed", 45000, "2026-03-09T11:20:00Z"},
	{"PAY-0005", "Sinar Elektronik", "processing", 1575000, "2026-03-08T17:00:00Z"},
	{"PAY-0006", "Pustaka Kita", "completed", 76000, "2026-03-08T10:10:00Z"},
	{"PAY-0007", "Ruang Hijau", "failed", 98000, "2026-03-07T16:25:00Z"},
	{"PAY-0008", "Dapur Ibu", "completed", 132000, "2026-03-07T09:40:00Z"},
	{"PAY-0009", "Lensa Studio", "processing", 325000, "2026-03-06T14:30:00Z"},
	{"PAY-0010", "Jaya Motor", "completed", 850000, "2026-03-06T08:05:00Z"},
	{"PAY-0011", "Teras Kayu", "failed", 65000, "2026-03-05T18:00:00Z"},
	{"PAY-0012", "Cahaya Optik", "completed", 275000, "2026-03-05T13:15:00Z"},
	{"PAY-0013", "Bumi Organik", "processing", 119000, "2026-03-04T15:50:00Z"},
	{"PAY-0014", "Sari Roti", "completed", 54000, "2026-03-04T07:35:00Z"},
	{"PAY-0015", "Nusa Apparel", "failed", 430000, "2026-03-03T19:10:00Z"},
	{"PAY-0016", "Kreasi Digital", "completed", 720000, "2026-03-03T10:25:00Z"},
	{"PAY-0017", "Pasar Segar", "processing", 67000, "2026-03-02T16:40:00Z"},
	{"PAY-0018", "Mitra Bangunan", "failed", 2300000, "2026-03-02T08:30:00Z"},
	{"PAY-0019", "Titik Temu", "completed", 185000, "2026-03-01T17:15:00Z"},
	{"PAY-0020", "Alam Lestari", "processing", 94000, "2026-03-01T09:00:00Z"},
	{"PAY-0021", "Raja Buah", "failed", 128000, "2026-02-28T18:45:00Z"},
	{"PAY-0022", "Satu Atap", "completed", 630000, "2026-02-28T13:20:00Z"},
	{"PAY-0023", "Kanvas Ruang", "processing", 360000, "2026-02-27T11:10:00Z"},
	{"PAY-0024", "Bersih Rumah", "completed", 82000, "2026-02-27T07:50:00Z"},
	{"PAY-0025", "Pelita Buku", "failed", 147000, "2026-02-26T15:35:00Z"},
	{"PAY-0026", "Kedai Senja", "completed", 96000, "2026-02-26T10:05:00Z"},
	{"PAY-0027", "Maju Bersama", "processing", 1125000, "2026-02-25T19:30:00Z"},
	{"PAY-0028", "Bunga Pagi", "failed", 53000, "2026-02-25T12:00:00Z"},
	{"PAY-0029", "Roda Dua", "completed", 450000, "2026-02-24T14:55:00Z"},
	{"PAY-0030", "Rasa Nusantara", "processing", 158000, "2026-02-24T08:20:00Z"},
}
