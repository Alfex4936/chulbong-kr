package kakao

// KakaoResponse is the top-level structure
type KakaoResponse struct {
	Meta      Meta       `json:"meta"`
	Documents []Document `json:"documents"`
}

// KakaoRegionResponse is the top-level structure
type KakaoRegionResponse struct {
	Meta      Meta          `json:"meta"`
	Documents []GeoDocument `json:"documents"`
}

// Meta contains response related info
type Meta struct {
	TotalCount    int   `json:"total_count"`
	PageableCount *int  `json:"pageable_count"`
	IsEnd         *bool `json:"is_end"`
}

// Document contains address details
type Document struct {
	Address     *Address     `json:"address"`
	RoadAddress *RoadAddress `json:"road_address"`

	AddressName *string `json:"address_name"`
	AddressType *string `json:"address_type"`
	X           *string `json:"x"`
	Y           *string `json:"y"`
}

// Address contains detailed info about the 지번 address
type Address struct {
	AddressName      string `json:"address_name"`
	Region1depthName string `json:"region_1depth_name"`
	Region2depthName string `json:"region_2depth_name"`
	Region3depthName string `json:"region_3depth_name"`
	MountainYN       string `json:"mountain_yn"`
	MainAddressNo    string `json:"main_address_no"`
	SubAddressNo     string `json:"sub_address_no"`
	ZipCode          string `json:"zip_code"` // Deprecated
}

// RoadAddress contains detailed info about the 도로명 address
type RoadAddress struct {
	AddressName      string `json:"address_name"`
	Region1depthName string `json:"region_1depth_name"`
	Region2depthName string `json:"region_2depth_name"`
	Region3depthName string `json:"region_3depth_name"`
	RoadName         string `json:"road_name"`
	UndergroundYN    string `json:"underground_yn"`
	MainBuildingNo   string `json:"main_building_no"`
	SubBuildingNo    string `json:"sub_building_no"`
	BuildingName     string `json:"building_name"`
	ZoneNo           string `json:"zone_no"`
}

// GeoDocument represents the structure of a geographical location document.
type GeoDocument struct {
	RegionType       string  `json:"region_type"`        // H(행정동) 또는 B(법정동)
	AddressName      string  `json:"address_name"`       // 전체 지역 명칭
	Region1DepthName string  `json:"region_1depth_name"` // 지역 1Depth, 시도 단위. 바다 영역은 존재하지 않음
	Region2DepthName string  `json:"region_2depth_name"` // 지역 2Depth, 구 단위. 바다 영역은 존재하지 않음
	Region3DepthName string  `json:"region_3depth_name"` // 지역 3Depth, 동 단위. 바다 영역은 존재하지 않음
	Region4DepthName string  `json:"region_4depth_name"` // 지역 4Depth, region_type이 법정동이며, 리 영역인 경우만 존재
	Code             string  `json:"code"`               // region 코드
	X                float64 `json:"x"`                  // X 좌표값, 경위도인 경우 경도(longitude)
	Y                float64 `json:"y"`                  // Y 좌표값, 경위도인 경우 위도(latitude)
}
