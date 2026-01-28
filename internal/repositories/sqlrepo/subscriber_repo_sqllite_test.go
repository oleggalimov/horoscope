package sqlrepo

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"horoscope/internal/horoscope/model"
	"log"
	"os"
	"testing"
	"time"
)

var (
	repo *SubscriberRepository
)

func TestMain(m *testing.M) {
	//setup
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}
	migrateDb(db)
	repo = NewSubscriberRepo(db)

	//test
	code := m.Run()

	// teardown
	db.Close()
	os.Remove("test.db")
	os.Exit(code)
}

func migrateDb(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS subscribers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT NOT NULL,
		address TEXT NOT NULL,
		sign TEXT NOT NULL,
		created_at DATETIME NOT NULL
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("failed to migrate DB: %v", err)
	}
}

func TestSubscriberRepository_Add(t *testing.T) {
	//GIVEN
	subscriber := model.Subscriber{
		Type:      model.Telegram,
		Address:   "test_1",
		Sign:      model.ZodiacSign("q1q"),
		CreatedAt: time.Now(),
	}

	//WHEN
	err := repo.Add(context.Background(), &subscriber)

	//THEN
	if err != nil {
		t.Fatalf("failed to add subscriber: %v", err)
	}
	_, err = repo.FindByID(context.Background(), subscriber.ID)
	if err != nil {
		t.Fatalf("failed to find subscriber: %v", err)
	}

}
