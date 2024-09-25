package dto

import (
	"time"
)

type MarkerReportRequest struct {
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	NewLatitude  float64 `json:"newLatitude,omitempty"`
	NewLongitude float64 `json:"newLongitude,omitempty"`
	MarkerID     int     `json:"markerId"`
	UserID       int     `json:"userId"`
	Description  string  `json:"description"`
	DoesExist    bool    `json:"doesExist,omitempty"`
}

type MarkerReportResponse struct {
	Latitude     float64   `json:"latitude" db:"Latitude"`
	Longitude    float64   `json:"longitude" db:"Longitude"`
	NewLatitude  float64   `json:"newLatitude,omitempty" db:"NewLatitude"`
	NewLongitude float64   `json:"newLongitude,omitempty" db:"NewLongitude"`
	CreatedAt    time.Time `json:"createdAt" db:"CreatedAt"`
	ReportID     int       `json:"reportId" db:"ReportID"`
	MarkerID     int       `json:"markerId" db:"MarkerID"`
	UserID       *int      `json:"userId,omitempty" db:"UserID"` // Pointer to handle nullable UserID
	Description  string    `json:"description" db:"Description"`
	PhotoURLs    []string  `json:"photoUrls,omitempty"` // Array to store multiple photo URLs
	Status       string    `json:"status" db:"Status"`
	Address      string    `json:"address,omitempty" db:"Address"`
	DoesExist    bool      `json:"doesExist,omitempty" db:"DoesExist"`
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

// GroupedReportsResponse represents the response structure for grouped reports by MarkerID
// type GroupedReportsResponse struct {
// 	TotalReports int                    `json:"totalReports"`
// 	Markers      *orderedmap.OrderedMap `json:"markers"`
// }

// ReportWithPhotos is a data transfer object for reports including photos
type ReportWithPhotos struct {
	NewLatitude  float64   `json:"newLatitude,omitempty"`
	NewLongitude float64   `json:"newLongitude,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	ReportID     int       `json:"reportID"`
	Description  string    `json:"description"`
	Status       string    `json:"status"`
	Photos       []string  `json:"photos"`
	Address      string    `json:"address,omitempty"`
}

type MarkerWithLatestReport struct {
	MarkerID   int
	LatestDate time.Time
}

type Location struct {
	Latitude  float64 `db:"Latitude"`
	Longitude float64 `db:"Longitude"`
}

type MarkerWithReports struct {
	MarkerID int                `json:"markerID"`
	Reports  []ReportWithPhotos `json:"reports"`
}

type GroupedReportsResponse struct {
	TotalReports int                 `json:"totalReports"`
	Markers      []MarkerWithReports `json:"markers"`
}
