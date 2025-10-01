package rating_redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/098765432m/grpc-kafka/rating/internal/domain"
	"github.com/redis/go-redis/v9"
)

type RedisRatingCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisRatingCache(client *redis.Client) *RedisRatingCache {
	return &RedisRatingCache{
		client: client,
	}
}

func (rrc *RedisRatingCache) GetRatingsByHotelId(ctx context.Context, hotelId string) ([]domain.Rating, error) {
	key := fmt.Sprintf("rating:hotel:%s", hotelId)
	val, err := rrc.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var ratings []domain.Rating
	if err := json.Unmarshal([]byte(val), &ratings); err != nil {
		return nil, err
	}

	return ratings, nil
}

func (rrs *RedisRatingCache) SetRatingsByHotelId(ctx context.Context, hotelId string, ratings []domain.Rating) error {
	key := fmt.Sprintf("rating:hotel:%s", hotelId)

	data, err := json.Marshal(ratings)
	if err != nil {
		return err
	}

	return rrs.client.Set(ctx, key, data, rrs.ttl).Err()
}
