package redis

import (
	"2023_2_Holi/domain"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type recomRedisRepository struct {
	client *redis.Client
}

func NewRecomRedisRepository(client *redis.Client) domain.RecomRepository {
	return &recomRedisRepository{client}
}

func (r *recomRedisRepository) GetRecommendations(userID int) ([]int, error) {
	key := fmt.Sprintf("recommendations:%d", userID)
	recommendations, err := r.client.LRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	// Преобразование списка строк в слайс целых чисел
	var intRecommendations []int
	for _, rec := range recommendations {
		movieId, err := strconv.Atoi(rec)
		if err != nil {
			return nil, err
		}
		intRecommendations = append(intRecommendations, movieId)
	}
	return intRecommendations, nil
}
