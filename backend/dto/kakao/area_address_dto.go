package kakao

// Coordinate represents a single coordinate point.
type Coordinate struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Shape defines the polygonal shape of a region.
type Shape struct {
	Hole           bool           `json:"hole"`
	Type           string         `json:"type"`
	CoordinateList [][]Coordinate `json:"coordinateList"`
}

// AreaDocument represents the document structure with optional and existing fields.
type AreaDocument struct {
	DocID    *string `json:"docid,omitempty"`
	Name     *string `json:"name,omitempty"`
	RoadName *string `json:"roadName,omitempty"`
	Building *string `json:"building,omitempty"`
	Bunji    *string `json:"bunji,omitempty"`
	Ho       *string `json:"ho,omitempty"`
	San      *string `json:"san,omitempty"`
	Shape    Shape   `json:"shape"`
	ZoneNo   *string `json:"zone_no,omitempty"`
}

// KakaoMarkerData holds all the information about a marker, including old and new document data, and other properties.
type KakaoMarkerData struct {
	Old           *AreaDocument `json:"old,omitempty"`
	X             float64       `json:"x"`
	Y             float64       `json:"y"`
	RegionID      string        `json:"regionid"`
	Region        string        `json:"region"`
	BCode         string        `json:"bcode"`
	LinePtsFormat string        `json:"line_pts_format"`
	New           *AreaDocument `json:"new,omitempty"`
}
