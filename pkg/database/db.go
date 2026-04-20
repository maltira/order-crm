package database

import (
	"database/sql"
	"log"
	"order-crm/config"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("pgx", config.Env.DatabaseUrl)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Проверка соединения
	if err = db.Ping(); err != nil {
		log.Fatal("Cannot ping database:", err)
	}
}

func GetDB() *sql.DB {
	return db
}
