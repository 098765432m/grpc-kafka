package rating_handler

import (
	"context"

	"github.com/098765432m/grpc-kafka/common/gen-proto/rating_pb"
	rating_repo "github.com/098765432m/grpc-kafka/rating/internal/repository/rating"
	rating_service "github.com/098765432m/grpc-kafka/rating/internal/service"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (rg *RatingGrpcHandler) GetRatingsByHotelId(ctx context.Context, req *rating_pb.GetRatingsByHotelIdRequest) (*rating_pb.GetRatingsByHotelIdRepsonse, error) {
	var hotelId pgtype.UUID
	if err := hotelId.Scan(req.GetHotelId()); err != nil {
		zap.S().Errorln("Failed to convert UUID")
		return nil, status.Error(codes.InvalidArgument, "Loi UUID khach san")
	}

	ratings, err := rg.service.GetRatingsByHotelId(ctx, hotelId)
	if err != nil {
		zap.S().Errorln(err)
		return nil, status.Error(codes.Internal, "Loi he thong")
	}

	var grpc_ratings []*rating_pb.Rating
	for _, rating := range ratings {
		grpc_rating := &rating_pb.Rating{
			Id:      rating.ID.String(),
			Score:   rating.Score,
			HotelId: rating.HotelID.String(),
			UserId:  rating.UserID.String(),
			Comment: rating.Comment.String,
		}
		grpc_ratings = append(grpc_ratings, grpc_rating)
	}

	return &rating_pb.GetRatingsByHotelIdRepsonse{
		Ratings: grpc_ratings,
	}, nil
}

func (rg *RatingGrpcHandler) UpdateRatingById(ctx context.Context, req *rating_pb.UpdateRatingByIdRequest) (*rating_pb.UpdateRatingByIdResponse, error) {

	err := rg.service.UpdateRating(ctx, &rating_repo.UpdateRatingParams{
		ID:      req.NewRating.Id,
		Score:   req.NewRating.Score,
		HotelID: req.NewRating.HotelId,
		UserID:  req.NewRating.UserId,
		Comment: req.NewRating.Comment,
	})

	if err != nil {
		zap.S().Errorln(err)
		return nil, err
	}

	return &rating_pb.UpdateRatingByIdResponse{}, nil
}

func (rg *RatingGrpcHandler) DeleteRatingById(ctx context.Context, req *rating_pb.DeleteRatingByIdRequest) (*rating_pb.DeleteRatingByIdResponse, error) {

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

	return &rating_pb.DeleteRatingByIdResponse{}, nil
}
