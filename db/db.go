package db

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Connect() {
	var err error
	// ساخت فایل دیتابیس اگر وجود نداشت
	if _, err := os.Stat("cafe.db"); os.IsNotExist(err) {
		file, _ := os.Create("cafe.db")
		file.Close()
	}
	DB, err = sql.Open("sqlite", "cafe.db")
	if err != nil {
		log.Fatal("خطا در اتصال به پایگاه داده:", err)
	}
}
