package sqlrepo

import (
	"context"
	"database/sql"
	"fmt"
	"horoscope/internal/horoscope/model"
	"log"
	"strings"
)

type SubscriberRepository struct {
	db *sql.DB
}

func NewSubscriberRepo(db *sql.DB) *SubscriberRepository {
	return &SubscriberRepository{db: db}
}

func (sr *SubscriberRepository) Add(ctx context.Context, s *model.Subscriber) error {
	query := `
        INSERT INTO subscribers (type, address, channel_id, sign, created_at, status)
        VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
    `
	row := sr.db.QueryRowContext(ctx, query, s.Type, s.Address, s.ChannelId, s.Sign, s.CreatedAt, s.Status)
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
		if err := rows.Scan(&s.ID, &s.Type, &s.Address, &s.ChannelId, &s.Sign, &s.CreatedAt, &s.Status); err != nil {
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
		`SELECT * FROM subscribers WHERE id = $1`, id).Scan(&s.ID, &s.Type, &s.Address, &s.ChannelId,
		&s.Sign, &s.CreatedAt, &s.Status)
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
		if err := rows.Scan(&s.ID, &s.Type, &s.Address, &s.ChannelId, &s.Sign, &s.CreatedAt, s.Status); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}
	return subs, nil
}

func (sr *SubscriberRepository) BatchUpdateAll(ctx context.Context, subscribers []model.Subscriber) (int, error) {

	if len(subscribers) == 0 {
		return 0, nil
	}

	tx, err := sr.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	args := make([]any, 0, len(subscribers)*6)
	values := make([]string, 0, len(subscribers))

	for _, s := range subscribers {
		values = append(values, "(?, ?, ?, ?, ?, ?)")
		args = append(args,
			s.Type,
			s.Address,
			s.ChannelId,
			s.Sign,
			s.CreatedAt,
			s.Status,
		)
	}

	query := fmt.Sprintf(`
	INSERT INTO subscribers (
		type,
		address,
		sign,
        channel_id,
		created_at,
		status
	)
	VALUES %s
	ON CONFLICT(type, address) DO UPDATE SET
		sign   = excluded.sign,
		status = excluded.status
	WHERE
		subscribers.sign   IS NOT excluded.sign OR
		subscribers.status IS NOT excluded.status
	`, strings.Join(values, ","))

	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return int(affected), nil
}
