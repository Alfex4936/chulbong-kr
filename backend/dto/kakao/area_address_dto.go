package kakao

// Coordinate represents a single coordinate point.
type Coordinate struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Shape defines the polygonal shape of a region.
type Shape struct {
	CoordinateList [][]Coordinate `json:"coordinateList"`
	Type           string         `json:"type"`
	Hole           bool           `json:"hole"`
}

// AreaDocument represents the document structure with optional and existing fields.
type AreaDocument struct {
	Shape    Shape   `json:"shape"`
	DocID    *string `json:"docid,omitempty"`
	Name     *string `json:"name,omitempty"`
	RoadName *string `json:"roadName,omitempty"`
	Building *string `json:"building,omitempty"`
	Bunji    *string `json:"bunji,omitempty"`
	Ho       *string `json:"ho,omitempty"`
	San      *string `json:"san,omitempty"`
	ZoneNo   *string `json:"zone_no,omitempty"`
}

// KakaoMarkerData holds all the information about a marker, including old and new document data, and other properties.
type KakaoMarkerData struct {
	X             float64       `json:"x"`
	Y             float64       `json:"y"`
	Old           *AreaDocument `json:"old,omitempty"`
	New           *AreaDocument `json:"new,omitempty"`
	RegionID      string        `json:"regionid"`
	Region        string        `json:"region"`
	BCode         string        `json:"bcode"`
	LinePtsFormat string        `json:"line_pts_format"`
}
