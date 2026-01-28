package sqlrepo

import (
	"context"
	"database/sql"
	"fmt"
	"horoscope/internal/horoscope/model"
	"log"
)

type SubscriberRepository struct {
	db *sql.DB
}

func NewSubscriberRepo(db *sql.DB) *SubscriberRepository {
	return &SubscriberRepository{db: db}
}

func (sr *SubscriberRepository) Add(ctx context.Context, s *model.Subscriber) error {
	query := `
        INSERT INTO subscribers (type, address, sign, created_at)
        VALUES ($1, $2, $3, $4) RETURNING id
    `
	row := sr.db.QueryRowContext(ctx, query, s.Type, s.Address, s.Sign, s.CreatedAt)
	if err := row.Scan(&s.ID); err != nil {
		return err
	}
	return nil
}

func (sr *SubscriberRepository) Remove(ctx context.Context, id int64) (int, error) {
	query := `
		DELETE FROM subscribers WHERE id = $1 RETURNING id
	`

	res, err := sr.db.ExecContext(ctx, query, id)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	if rows == 0 {
		return 0, fmt.Errorf("подписчик не найден")
	}
	return int(rows), nil
}

func (sr *SubscriberRepository) FindAll(ctx context.Context) ([]model.Subscriber, error) {
	rows, err := sr.db.QueryContext(ctx, `SELECT * FROM subscribers`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var subs []model.Subscriber
	for rows.Next() {
		var s model.Subscriber
		if err := rows.Scan(&s.ID, &s.Type, &s.Address, &s.Sign, &s.CreatedAt); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}
	return subs, nil
}

func (sr *SubscriberRepository) FindByID(ctx context.Context, id int64) (model.Subscriber, error) {
	var s model.Subscriber
	err := sr.db.QueryRowContext(
		ctx,
		`SELECT * FROM subscribers WHERE id = $1`, id).Scan(&s.ID, &s.Type, &s.Address, &s.Sign, &s.CreatedAt)
	if err != nil {
		return model.Subscriber{}, err
	}
	return s, nil
}

func (sr *SubscriberRepository) FindByType(_ context.Context, t model.Channel) ([]model.Subscriber, error) {
	rows, err := sr.db.Query("SELECT * FROM subscribers WHERE type = $1", t)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal("Ошибка при чтении списка", err)
		}
	}(rows)
	var subs []model.Subscriber
	for rows.Next() {
		var s model.Subscriber
		if err := rows.Scan(&s.ID, &s.Type, &s.Address, &s.Sign, &s.CreatedAt); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}
	return subs, nil
}
