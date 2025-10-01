package rating_service

import (
	"context"

	rating_domain "github.com/098765432m/grpc-kafka/rating/internal/domain"
	rating_redis "github.com/098765432m/grpc-kafka/rating/internal/infrastructure/redis"
	rating_repo "github.com/098765432m/grpc-kafka/rating/internal/infrastructure/repository/sqlc/rating"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type RatingService struct {
	repo  *rating_repo.Queries
	redis *rating_redis.RedisRatingCache
}

func NewRatingService(repo *rating_repo.Queries, redis *rating_redis.RedisRatingCache) *RatingService {
	return &RatingService{
		repo:  repo,
		redis: redis,
	}
}

func (rs *RatingService) GetRatingsByHotelId(ctx context.Context, hotelId pgtype.UUID) ([]rating_domain.Rating, error) {

	cache, err := rs.redis.GetRatingsByHotelId(ctx, hotelId.String())
	if err != nil {
		zap.S().Errorln("Failed to get Redis Ratings by Hotel id: ", err)
		return nil, err
	}

	if cache != nil {
		return cache, nil
	}

	ratings, err := rs.repo.GetRatingsByHotel(ctx, hotelId)
	if err != nil {
		zap.S().Errorln("Failed to get Ratings by Hotel id: ", err)
		return nil, err
	}

	result := make([]rating_domain.Rating, 0, len(ratings))
	for _, rating := range ratings {
		result = append(result, rating_domain.Rating{
			Id:      rating.ID.String(),
			Score:   int(rating.Score),
			HotelId: rating.HotelID.String(),
			UserId:  rating.UserID.String(),
			Comment: rating.Comment.String,
		})

	}

	return result, nil
}

func (rs *RatingService) CreateRating(ctx context.Context, newRating *rating_repo.CreateRatingParams) error {

	err := rs.repo.CreateRating(ctx, *newRating)
	if err != nil {
		zap.S().Errorln("Failed to create new Rating")
		return err
	}

	return nil
}

func (rs *RatingService) UpdateRating(ctx context.Context, updateRatingParams *rating_repo.UpdateRatingParams) error {

	err := rs.repo.UpdateRating(ctx, *updateRatingParams)
	if err != nil {
		zap.S().Errorln("Failed to update Rating by id")
		return err
	}

	return nil
}

func (rs *RatingService) DeleteRatingById(ctx context.Context, id pgtype.UUID) error {

	err := rs.repo.DeleteRating(ctx, id)
	if err != nil {
		zap.S().Errorln("Failed to Delete Rating by id")
		return err
	}
	return nil
}
