package redis

import (
	"2023_2_Holi/domain"
	"github.com/redis/go-redis/v9"
)

type recomRedisRepository struct {
	client *redis.Client
}

func NewRecomRedisRepository(client *redis.Client) domain.RecomRepository {
	return &recomRedisRepository{client}
}

func (r *recomRedisRepository) GetRecommendations(userID int) ([]domain.Video, error) {

}
