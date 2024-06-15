package dto

import "time"

type MarkerReportRequest struct {
	MarkerID     int     `json:"markerId"`
	UserID       int     `json:"userId"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	NewLatitude  float64 `json:"newLatitude,omitempty"`
	NewLongitude float64 `json:"newLongitude,omitempty"`
	Description  string  `json:"description"`
}

type MarkerReportResponse struct {
	ReportID     int       `json:"reportId" db:"ReportID"`
	MarkerID     int       `json:"markerId" db:"MarkerID"`
	UserID       *int      `json:"userId,omitempty" db:"UserID"` // Pointer to handle nullable UserID
	Latitude     float64   `json:"latitude" db:"Latitude"`
	Longitude    float64   `json:"longitude" db:"Longitude"`
	NewLatitude  float64   `json:"newLatitude,omitempty" db:"NewLatitude"`
	NewLongitude float64   `json:"newLongitude,omitempty" db:"NewLongitude"`
	Description  string    `json:"description" db:"Description"`
	PhotoURLs    []string  `json:"photoUrls,omitempty"` // Array to store multiple photo URLs
	CreatedAt    time.Time `json:"createdAt" db:"CreatedAt"`
	Status       string    `json:"status" db:"Status"`
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
type GroupedReportsResponse struct {
	TotalReports int                        `json:"totalReports"`
	Markers      map[int][]ReportWithPhotos `json:"markers"`
}

// ReportWithPhotos is a data transfer object for reports including photos
type ReportWithPhotos struct {
	ReportID    int       `json:"reportID"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	Photos      []string  `json:"photos"`
}
