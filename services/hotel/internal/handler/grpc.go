package handler

import (
	"context"

	"github.com/098765432m/grpc-kafka/common/gen-proto/hotel_pb"
	"github.com/098765432m/grpc-kafka/hotel/internal/service"
)

type HotelGrpcHandler struct {
	hotel_pb.UnimplementedHotelServiceServer
	service *service.HotelService
}

func NewHotelGrpcHandler(service *service.HotelService) *HotelGrpcHandler {
	return &HotelGrpcHandler{
		service: service,
	}
}

func (hg *HotelGrpcHandler) GetHotel(ctx context.Context, req *hotel_pb.GetHotelRequest) (*hotel_pb.GetHotelResponse, error) {

	return &hotel_pb.GetHotelResponse{
		Id:   "123",
		Name: "Sample Hotel",
	}, nil
}

func (hg *HotelGrpcHandler) CreateHotel(ctx context.Context, req *hotel_pb.CreateHotelRequest) (*hotel_pb.CreateHotelResponse, error) {
	err := hg.service.CreateHotel(ctx, req.Name)
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

	var grpc_Hotels []*hotel_pb.Hotel
	for _, hotel := range hotels {
		grpc_hotel := &hotel_pb.Hotel{
			Id:   hotel.ID.String(),
			Name: hotel.Name,
		}

		grpc_Hotels = append(grpc_Hotels, grpc_hotel)
	}

	return &hotel_pb.GetAllHotelsResponse{
		Hotels: grpc_Hotels,
	}, nil
}
