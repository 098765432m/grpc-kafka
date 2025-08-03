package rating_service

import (
	"context"

	rating_repo "github.com/098765432m/grpc-kafka/rating/internal/repository/rating"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type RatingService struct {
	repo *rating_repo.Queries
}

func NewRatingService(repo *rating_repo.Queries) *RatingService {
	return &RatingService{
		repo: repo,
	}
}

func (rs *RatingService) GetRatingsByHotel(ctx context.Context, hotelId pgtype.UUID) ([]rating_repo.Rating, error) {

	ratings, err := rs.repo.GetRatingsByHotel(ctx, hotelId)
	if err != nil {
		zap.S().Errorln("Failed to get Ratings by Hotel id")
		return nil, err
	}

	return ratings, nil
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
