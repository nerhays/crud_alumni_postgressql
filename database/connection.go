package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	var err error

	// Ambil langsung dari .env
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:12345@localhost:5432/alumni_db?sslmode=disable"
	}

	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Gagal koneksi DB:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("DB tidak bisa di-ping:", err)
	}

	log.Println("✅ Berhasil konek ke PostgreSQL")
}



