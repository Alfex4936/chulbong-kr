package dto

import "time"

type MarkerReportRequest struct {
	MarkerID    int     `json:"markerId"`
	UserID      int     `json:"userId"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Description string  `json:"description"`
}

type MarkerReportResponse struct {
	ReportID       int       `json:"-" db:"ReportID"`
	MarkerID       int       `json:"markerId" db:"MarkerID"`
	UserID         *int      `json:"userId,omitempty" db:"UserID"` // Pointer to handle nullable UserID
	Latitude       float64   `json:"latitude" db:"Latitude" `
	Longitude      float64   `json:"longitude" db:"Longitude"`
	Description    string    `json:"description" db:"Description"`
	ReportImageURL string    `json:"reportImageUrl,omitempty" db:"ReportImageURL"`
	CreatedAt      time.Time `json:"createdAt" db:"CreatedAt"`
}

// MarkerReports groups all reports for a specific marker.
type MarkerReports struct {
	Reports []MarkerReportResponse `json:"reports"`
}

// ReportsResponse is the structured response for all reports.
type ReportsResponse struct {
	TotalReports int                   `json:"totalReports"`
	Markers      map[int]MarkerReports `json:"markers"`
}

// // Convert your latitude and longitude into a geom.Point for easier handling of geographic data.
// func NewPoint(lat, lon float64) *geom.Point {
// 	return geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{X: lon, Y: lat})
// }
