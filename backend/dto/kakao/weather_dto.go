package kakao

// WeatherResponse
type WeatherResponse struct {
	Codes        Codes        `json:"codes"`
	WeatherInfos WeatherInfos `json:"weatherInfos"`
}

type Codes struct {
	Hcode      Code   `json:"hcode"`
	Bcode      Code   `json:"bcode"`
	ResultCode string `json:"resultCode"`
}

type Code struct {
	Childcount float64 `json:"childcount"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	Type       string  `json:"type"`
	Code       string  `json:"code"`
	Name       string  `json:"name"`
	FullName   string  `json:"fullName"`
	RegionId   string  `json:"regionId"`
	Name0      string  `json:"name0"`
	Code1      string  `json:"code1"`
	Name1      string  `json:"name1"`
	Code2      string  `json:"code2"`
	Name2      string  `json:"name2"`
	Code3      string  `json:"code3"`
	Name3      string  `json:"name3"`
}

type WeatherInfos struct {
	Current  WeatherInfo `json:"current"`
	Forecast WeatherInfo `json:"forecast"`
}

type WeatherInfo struct {
	Type        string `json:"type"`
	Rcode       string `json:"rcode"`
	IconId      string `json:"iconId"`
	Temperature string `json:"temperature"`
	Desc        string `json:"desc"`
	Humidity    string `json:"humidity"`
	Rainfall    string `json:"rainfall"`
	Snowfall    string `json:"snowfall"`
}

type WeatherRequest struct {
	Temperature string `json:"temperature"`
	Desc        string `json:"desc"`
	IconImage   string `json:"iconImage"`
	Humidity    string `json:"humidity"`
	Rainfall    string `json:"rainfall"`
	Snowfall    string `json:"snowfall"`
}
