package redis

import (
	"2023_2_Holi/domain"
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type utilsRedisRepository struct {
	client *redis.Client
}

func NewUtilsRedisRepository(c *redis.Client) domain.UtilsRepository {
	return &utilsRedisRepository{
		client: c,
	}
}

func (s *utilsRedisRepository) GetIdBy(token string) (int, error) {
	if token == "" {
		return 0, domain.ErrInvalidToken
	}

	r, err := s.client.Get(context.Background(), token).Result()
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(r)
}
