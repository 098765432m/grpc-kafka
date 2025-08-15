package hotel_handler

import (
	"context"

	"github.com/098765432m/grpc-kafka/common/gen-proto/hotel_pb"
	hotel_repo "github.com/098765432m/grpc-kafka/hotel/internal/repository/hotel"
	hotel_service "github.com/098765432m/grpc-kafka/hotel/internal/service/hotel"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type HotelGrpcHandler struct {
	hotel_pb.UnimplementedHotelServiceServer
	service *hotel_service.HotelService
}

func NewHotelGrpcHandler(service *hotel_service.HotelService) *HotelGrpcHandler {
	return &HotelGrpcHandler{
		service: service,
	}
}

func (hg *HotelGrpcHandler) GetHotelById(ctx context.Context, req *hotel_pb.GetHotelByIdRequest) (*hotel_pb.GetHotelByIdResponse, error) {

	var id pgtype.UUID
	if err := id.Scan(req.Id); err != nil {
		zap.S().Errorln("Failed to Hotel UUID")
		return nil, err
	}

	hotel, err := hg.service.GetHotelById(ctx, id)
	if err != nil {

		zap.S().Errorln("Failed to get Hotel by Id: ", err)
		return nil, err
	}

	return &hotel_pb.GetHotelByIdResponse{
		Hotel: &hotel_pb.Hotel{
			Id:      hotel.ID.String(),
			Name:    hotel.Name,
			Address: hotel.Address,
		},
	}, nil
}

func (hg *HotelGrpcHandler) CreateHotel(ctx context.Context, req *hotel_pb.CreateHotelRequest) (*hotel_pb.CreateHotelResponse, error) {
	err := hg.service.CreateHotel(ctx, &hotel_repo.CreateHotelParams{
		Name:    req.Name,
		Address: req.Address,
	})
	if err != nil {
		return nil, err
	}
	return &hotel_pb.CreateHotelResponse{}, nil
}

func (hg *HotelGrpcHandler) GetAllHotels(ctx context.Context, req *hotel_pb.GetAllHotelsRequest) (*hotel_pb.GetAllHotelsResponse, error) {
	hotels, err := hg.service.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var grpc_hotels []*hotel_pb.Hotel
	for _, hotel := range hotels {
		grpc_hotel := &hotel_pb.Hotel{
			Id:   hotel.ID.String(),
			Name: hotel.Name,
		}

		grpc_hotels = append(grpc_hotels, grpc_hotel)
	}

	return &hotel_pb.GetAllHotelsResponse{
		Hotels: grpc_hotels,
	}, nil
}
