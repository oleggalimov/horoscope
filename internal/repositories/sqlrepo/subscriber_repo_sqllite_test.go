package sqlrepo

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"horoscope/internal/config"
	"horoscope/internal/horoscope/model"
	"log"
	"os"
	"testing"
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
	config.MigrateDb(db)
	repo = NewSubscriberRepo(db)

	//test
	code := m.Run()

	// teardown
	db.Close()
	os.Remove("test.db")
	os.Exit(code)
}

func TestSubscriberRepository_Add(t *testing.T) {
	//GIVEN
	subscriber := model.Subscriber{
		Type:      model.Telegram,
		Address:   "test_1",
		Sign:      model.ZodiacSign("q1q"),
		CreatedAt: nil,
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
