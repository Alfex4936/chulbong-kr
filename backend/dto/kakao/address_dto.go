package kakao

// KakaoResponse is the top-level structure
type KakaoResponse struct {
	Meta      Meta       `json:"meta"`
	Documents []Document `json:"documents"`
}

// Meta contains response related info
type Meta struct {
	TotalCount int `json:"total_count"`
}

// Document contains address details
type Document struct {
	Address     *Address     `json:"address"`
	RoadAddress *RoadAddress `json:"road_address"`
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
