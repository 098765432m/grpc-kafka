package hotel_repo_mapping

import (
	hotel_domain "github.com/098765432m/grpc-kafka/hotel/internal/domain"
	hotel_repo "github.com/098765432m/grpc-kafka/hotel/internal/infrastructure/repository/sqlc/hotel"
	room_repo "github.com/098765432m/grpc-kafka/hotel/internal/infrastructure/repository/sqlc/room"
	room_type_repo "github.com/098765432m/grpc-kafka/hotel/internal/infrastructure/repository/sqlc/room-type"
)

func FromHotelRepoToHotelDomain(hotelRepo hotel_repo.Hotel) hotel_domain.Hotel {

	return hotel_domain.Hotel{
		Id:      hotelRepo.ID.String(),
		Name:    hotelRepo.Name,
		Address: hotelRepo.Address.String,
	}
}

func FromHotelsRepoToHotelsDomain(hotelsRepo []hotel_repo.Hotel) []hotel_domain.Hotel {

	hotels := make([]hotel_domain.Hotel, 0, len(hotelsRepo))

	for _, h := range hotelsRepo {
		hotels = append(hotels, FromHotelRepoToHotelDomain(h))
	}

	return hotels
}

func FromRoomTypeRepoToRoomTypeDomain(roomTypeRepo room_type_repo.RoomType) hotel_domain.RoomType {
	return hotel_domain.RoomType{
		Id:      roomTypeRepo.ID.String(),
		Name:    roomTypeRepo.Name,
		Price:   int(roomTypeRepo.Price),
		HotelId: roomTypeRepo.HotelID.String(),
	}
}

func FromRoomTypesRepoToRoomTypesDomain(roomTypesRepo []room_type_repo.RoomType) []hotel_domain.RoomType {
	roomTypes := make([]hotel_domain.RoomType, 0, len(roomTypesRepo))

	for _, r := range roomTypesRepo {
		roomTypes = append(roomTypes, FromRoomTypeRepoToRoomTypeDomain(r))
	}

	return roomTypes
}

func FromRoomRepoToRoomDomain(roomRepo room_repo.Room) hotel_domain.Room {

	return hotel_domain.Room{
		Id:         roomRepo.ID.String(),
		Name:       roomRepo.Name,
		Status:     string(roomRepo.Status.RoomStatus),
		RoomTypeId: roomRepo.RoomTypeID.String(),
		HotelId:    roomRepo.HotelID.String(),
	}
}

func FromRoomsRepoToRoomsDomain(roomsRepo []room_repo.Room) []hotel_domain.Room {

	rooms := make([]hotel_domain.Room, 0, len(roomsRepo))
	for _, r := range roomsRepo {
		rooms = append(rooms, FromRoomRepoToRoomDomain(r))
	}

	return rooms
}
