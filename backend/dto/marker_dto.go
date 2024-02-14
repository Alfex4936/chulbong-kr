package dto

type MarkerRequest struct {
	MarkerID    int     `json:"markerId"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Description string  `json:"description"`
	PhotoURL    string  `json:"photoUrl"`
}

type MarkerResponse struct {
	MarkerID    int     `json:"markerId"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Description string  `json:"description"`
	Username    string  `json:"username"`
	PhotoURL    string  `json:"photoUrl"`
}
