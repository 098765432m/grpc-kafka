package rating_handler

import (
	"context"

	"github.com/098765432m/grpc-kafka/common/gen-proto/rating_pb"
	rating_repo "github.com/098765432m/grpc-kafka/rating/internal/repository/rating"
	rating_service "github.com/098765432m/grpc-kafka/rating/internal/service"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type RatingGrpcHandler struct {
	rating_pb.UnimplementedRatingServiceServer
	service *rating_service.RatingService
}

func NewRatingGrpcHandler(service *rating_service.RatingService) *RatingGrpcHandler {
	return &RatingGrpcHandler{
		service: service,
	}
}

func (rg *RatingGrpcHandler) GetRatingsByHotel(ctx context.Context, req *rating_pb.GetRatingsByHotelRequest) (*rating_pb.GetRatingsByHotelRepsonse, error) {
	var id pgtype.UUID
	if err := id.Scan(req.HotelId); err != nil {
		zap.S().Errorln("Failed to convert UUID")
		return nil, err
	}

	ratings, err := rg.service.GetRatingsByHotel(ctx, id)
	if err != nil {
		zap.S().Errorln(err)
		return nil, err
	}

	var grpc_ratings []*rating_pb.Rating
	for _, rating := range ratings {
		grpc_rating := &rating_pb.Rating{
			Id:      rating.ID.String(),
			Rating:  rating.Rating,
			HotelId: rating.HotelID.String(),
			UserId:  rating.UserID.String(),
			Comment: rating.Comment.String,
		}
		grpc_ratings = append(grpc_ratings, grpc_rating)
	}

	return &rating_pb.GetRatingsByHotelRepsonse{
		Ratings: grpc_ratings,
	}, nil
}

func (rg *RatingGrpcHandler) UpdateRating(ctx context.Context, req *rating_pb.UpdateRatingRequest) (*rating_pb.UpdateRatingResponse, error) {

	err := rg.service.UpdateRating(ctx, &rating_repo.UpdateRatingParams{
		ID:      req.NewRating.Id,
		Rating:  req.NewRating.Rating,
		HotelID: req.NewRating.HotelId,
		UserID:  req.NewRating.UserId,
		Comment: req.NewRating.Comment,
	})

	if err != nil {
		zap.S().Errorln(err)
		return nil, err
	}

	return &rating_pb.UpdateRatingResponse{}, nil
}

func (rg *RatingGrpcHandler) DeleteRating(ctx context.Context, req *rating_pb.DeleteRatingRequest) (*rating_pb.DeleteRatingResponse, error) {

	var id pgtype.UUID
	if err := id.Scan(req.Id); err != nil {
		zap.S().Errorln("Failed to convert UUID")
		return nil, err

	}

	err := rg.service.DeleteRatingById(ctx, id)
	if err != nil {
		zap.S().Errorln(err)
		return nil, err
	}

	return &rating_pb.DeleteRatingResponse{}, nil
}
