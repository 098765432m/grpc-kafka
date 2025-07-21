package handler

import (
	"context"

	"github.com/098765432m/grpc-kafka/common/gen-proto/hotels"
)

type HotelGrpcHandler struct {
	hotels.UnimplementedHotelServiceServer
}

func NewHotelGrpcHandler() *HotelGrpcHandler {
	return &HotelGrpcHandler{}
}

func (hg *HotelGrpcHandler) GetHotel(ctx context.Context, req *hotels.GetHotelRequest) (*hotels.GetHotelResponse, error) {

	return &hotels.GetHotelResponse{
		Id:      "123",
		Name:    "Sample Hotel",
		Address: "123 Sample St, Sample City",
	}, nil
}
