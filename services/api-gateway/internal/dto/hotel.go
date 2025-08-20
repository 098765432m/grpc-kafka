package api_dto

type HotelImage struct {
	Id       string `json:"id"`
	PublicId string `json:"public_id"`
	Format   string `json:"format"`
}

type HotelResponse struct {
	Id      string       `json:"id"`
	Name    string       `json:"name"`
	Address string       `json:"address"`
	Images  []HotelImage `json:"images"`
}
