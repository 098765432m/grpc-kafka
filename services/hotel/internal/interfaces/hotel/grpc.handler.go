package hotel_handler

import (
	"context"
	"errors"

	common_error "github.com/098765432m/grpc-kafka/common/error"
	"github.com/098765432m/grpc-kafka/common/gen-proto/hotel_pb"
	"github.com/098765432m/grpc-kafka/common/utils"
	hotel_service "github.com/098765432m/grpc-kafka/hotel/internal/application/hotel"
	room_service "github.com/098765432m/grpc-kafka/hotel/internal/application/room"
	room_type_service "github.com/098765432m/grpc-kafka/hotel/internal/application/room-type"
	hotel_repo "github.com/098765432m/grpc-kafka/hotel/internal/infrastructure/repository/sqlc/hotel"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HotelGrpcHandler struct {
	hotel_pb.UnimplementedHotelServiceServer
	service         *hotel_service.HotelService
	roomTypeService *room_type_service.RoomTypeService
	roomService     *room_service.RoomService
}

func NewHotelGrpcHandler(
	service *hotel_service.HotelService,
	roomTypeService *room_type_service.RoomTypeService,
	roomService *room_service.RoomService,
) *HotelGrpcHandler {
	return &HotelGrpcHandler{
		service:         service,
		roomTypeService: roomTypeService,
		roomService:     roomService,
	}
}

func (hg *HotelGrpcHandler) GetAllHotels(ctx context.Context, req *hotel_pb.GetAllHotelsRequest) (*hotel_pb.GetAllHotelsResponse, error) {
	hotels, err := hg.service.GetAll(ctx)
	if err != nil {
		zap.S().Info("Failed to get all hotels: ", err)
		return nil, status.Error(codes.Internal, "Loi khong lay duoc danh sach khach san")
	}

	var grpc_hotels []*hotel_pb.Hotel
	for _, hotel := range hotels {
		grpc_hotel := &hotel_pb.Hotel{
			Id:   hotel.Id,
			Name: hotel.Name,
		}

		grpc_hotels = append(grpc_hotels, grpc_hotel)
	}

	return &hotel_pb.GetAllHotelsResponse{
		Hotels: grpc_hotels,
	}, nil
}

func (hg *HotelGrpcHandler) GetHotelsByAddress(ctx context.Context, req *hotel_pb.GetHotelsByAddressRequest) (*hotel_pb.GetHotelsByAddressResponse, error) {
	address, err := utils.ParsePgText(req.GetAddress())
	if err != nil {
		zap.S().Infoln("Address is an invalid format")
		return nil, status.Error(codes.InvalidArgument, "Dia chi khong hop le")
	}

	hotelName, err := utils.ParsePgText(req.GetHotelName())
	if err != nil {
		zap.S().Infoln("Hotel Name is an invalid format")
		return nil, status.Error(codes.InvalidArgument, "Ten khach san khong hop le")
	}

	hotelIds, err := hg.service.GetHotelsByAddress(ctx, address, hotelName)
	if err != nil {
		return nil, status.Error(codes.Internal, "Loi khong lay duoc danh sach khach san")
	}

	hotelIdsStr := utils.ToPgUuidString(hotelIds)

	return &hotel_pb.GetHotelsByAddressResponse{
		HotelIds: hotelIdsStr,
	}, nil
}

func (hg *HotelGrpcHandler) GetHotelById(ctx context.Context, req *hotel_pb.GetHotelByIdRequest) (*hotel_pb.GetHotelByIdResponse, error) {

	var id pgtype.UUID
	if err := id.Scan(req.Id); err != nil {
		zap.S().Errorln("Invalid Hotel UUID")
		return nil, status.Error(codes.InvalidArgument, "Loi UUID khach san")
	}

	hotel, err := hg.service.GetHotelById(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, common_error.ErrNoRows):
			return nil, status.Error(codes.NotFound, "Khach san khong ton tai")
		}

		zap.S().Errorln("Failed to get Hotel by Id: ", err)
		return nil, err
	}

	return &hotel_pb.GetHotelByIdResponse{
		Hotel: &hotel_pb.Hotel{
			Id:      hotel.Id,
			Name:    hotel.Name,
			Address: hotel.Address,
		},
	}, nil
}

func (hg *HotelGrpcHandler) CreateHotel(ctx context.Context, req *hotel_pb.CreateHotelRequest) (*hotel_pb.CreateHotelResponse, error) {
	address, err := utils.ParsePgText(req.GetAddress())
	if err != nil {
		zap.S().Infoln("Address is an invalid format")
		return nil, status.Error(codes.InvalidArgument, "Dia chi khong hop le")
	}

	hotelName, err := utils.ParsePgText(req.GetName())
	if err != nil {
		zap.S().Infoln("Hotel Name is an invalid format")
		return nil, status.Error(codes.InvalidArgument, "Ten khach san khong hop le")
	}

	err = hg.service.CreateHotel(ctx, &hotel_repo.CreateHotelParams{
		Name:    hotelName,
		Address: address,
	})
	if err != nil {
		return nil, err
	}
	return &hotel_pb.CreateHotelResponse{}, nil
}

func (hg *HotelGrpcHandler) FilterHotels(ctx context.Context, req *hotel_pb.FilterHotelsRequest) (*hotel_pb.FilterHotelsResponse, error) {

	roomTypeIds, err := utils.ToPgUuidArray(req.GetRoomTypeIds())
	if err != nil {
		zap.S().Infoln("Invalid Room Type UUIDs")
		return nil, status.Error(codes.InvalidArgument, "uuid room type khong hop le")
	}

	zap.S().Infoln("Min Price REQ: ", req.GetMinPrice())
	zap.S().Infoln("Max Price REQ: ", req.GetMaxPrice())

	minPrice := utils.ToPgInt4(int(req.GetMinPrice()))
	maxPrice := utils.ToPgInt4(int(req.GetMaxPrice()))

	zap.S().Infoln("Min Price: ", minPrice)
	zap.S().Infoln("Max Price: ", maxPrice)

	rows, err := hg.service.FilterHotels(ctx, roomTypeIds, minPrice, maxPrice)
	if err != nil {
		return nil, status.Error(codes.Internal, "Loi khong filter duoc Hotels")
	}

	results := make([]*hotel_pb.FilterHotelRow, 0, len(rows))
	for _, row := range rows {
		results = append(results, &hotel_pb.FilterHotelRow{
			HotelId:      row.ID.String(),
			HotelName:    row.Name,
			HotelAddress: row.Address.String,
			MinPrice:     row.MinPrice.(int32), // assertion interface{} to int32
		})
	}

	return &hotel_pb.FilterHotelsResponse{
		FilterHotelRows: results,
	}, nil

}
