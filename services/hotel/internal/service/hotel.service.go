package service

type HotelService struct {
}

func NewHotelService() *HotelService {
	return &HotelService{}
}

func (hs *HotelService) GetHotelById(id string) error {
	return nil
}
