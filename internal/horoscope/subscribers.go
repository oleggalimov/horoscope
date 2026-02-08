package horoscope

import (
	"context"
	"horoscope/internal/horoscope/model"
)

type SubscribersDao interface {
	Add(ctx context.Context, s *model.Subscriber) error
	Remove(ctx context.Context, id int64) (int, error)
	FindAll(ctx context.Context) ([]model.Subscriber, error)
	FindByID(ctx context.Context, id int64) (model.Subscriber, error)
	FindByType(ctx context.Context, t model.Channel) ([]model.Subscriber, error)
	BatchUpdateAll(ctx context.Context, subscribers []model.Subscriber) (int, error)
}
