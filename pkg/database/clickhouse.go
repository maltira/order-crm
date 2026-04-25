package database

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

var CHDB *sql.DB

func InitClickHouse() {
	d, err := sql.Open("clickhouse", os.Getenv("CLICKHOUSE_URL"))
	if err != nil {
		log.Fatal("ClickHouse connection error:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("ClickHouse ping error:", err)
	}

	CHDB = d

	log.Println("ClickHouse connected")
}

func LogUserEvent(userID int, action string) {
	query := "INSERT INTO lab_crm_db.user_events (id, user_id, action, ts) VALUES (generateUUIDv4(), ?, ?, ?)"

	_, err := CHDB.Exec(query,
		userID,
		action,
		time.Now(),
	)

	if err != nil {
		log.Println("clickhouse user event error:", err)
	}
}

func LogOrderStatus(orderID int, statusID int) {
	query := `
		INSERT INTO lab_crm_db.orders_log (id, order_id, status_id, ts)
		VALUES (generateUUIDv4(), ?, ?, ?)
	`

	_, err := CHDB.Exec(query,
		orderID,
		statusID,
		time.Now(),
	)

	if err != nil {
		log.Println("clickhouse order log error:", err)
	}
}
