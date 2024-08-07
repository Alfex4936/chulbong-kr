package kakao

type StreetView struct {
	StreetList []Street    `json:"streetList,omitempty"`
	Street     interface{} `json:"street,omitempty"`
	Cnt        int         `json:"cnt,omitempty"`
}

type Street struct {
	Wtmx     float64     `json:"wtmx,omitempty"`
	Wtmy     float64     `json:"wtmy,omitempty"`
	Wgsx     float64     `json:"wgsx,omitempty"`
	Wgsy     float64     `json:"wgsy,omitempty"`
	Wcongx   float64     `json:"wcongx,omitempty"`
	Wcongy   float64     `json:"wcongy,omitempty"`
	ID       int         `json:"id,omitempty"`
	Angle    string      `json:"angle,omitempty"`
	ImgPath  string      `json:"img_path,omitempty"`
	Addr     string      `json:"addr,omitempty"`
	StName   string      `json:"st_name,omitempty"`
	StType   string      `json:"st_type,omitempty"`
	ShotDate string      `json:"shot_date,omitempty"`
	ShotTool string      `json:"shot_tool,omitempty"`
	AreaType interface{} `json:"area_type,omitempty"`
	Spot     interface{} `json:"spot,omitempty"`
	Past     interface{} `json:"past,omitempty"`
}

type StreetViewData struct {
	StreetView StreetView `json:"street_view,omitempty"`
}
