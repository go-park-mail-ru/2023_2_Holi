package redis

import (
	"2023_2_Holi/domain"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type sessionRedisRepository struct {
	client *redis.Client
}

func NewSessionRedisRepository(client *redis.Client) domain.SessionRepository {
	return &sessionRedisRepository{client}
}

func (s *sessionRedisRepository) Add(session domain.Session) error {
	duration := session.ExpiresAt.Sub(time.Now())
	err := s.client.Set(context.Background(), session.Token, session.UserID, duration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *sessionRedisRepository) DeleteByToken(token string) error {
	err := s.client.Del(context.Background(), token).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *sessionRedisRepository) SessionExists(token string) (bool, error) {
	exists, err := s.client.Exists(context.Background(), token).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}
