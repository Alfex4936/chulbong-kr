package model

// Facility represents the structure for a single facility record.
type Facility struct {
	FacilityID int `db:"FacilityID" json:"facilityId"`
	MarkerID   int `db:"MarkerID" json:"markerId"`
	Quantity   int `db:"Quantity" json:"quantity"`
}
