package facade

import (
	"github.com/Alfex4936/chulbong-kr/dto"
	"github.com/Alfex4936/chulbong-kr/dto/kakao"
)

// Get
func (mfs *MarkerFacadeService) FindClosestNMarkersWithinDistance(lat, lng float64, distance, pageSize, offset int) ([]dto.MarkerWithDistanceAndPhoto, int, error) {
	return mfs.LocationService.FindClosestNMarkersWithinDistance(lat, lng, distance, pageSize, offset)
}
func (mfs *MarkerFacadeService) FindRankedMarkersInCurrentArea(lat, lng float64, distance, limit int) ([]dto.MarkerWithDistanceAndPhoto, error) {
	return mfs.LocationService.FindRankedMarkersInCurrentArea(lat, lng, distance, limit)
}

func (mfs *MarkerFacadeService) FetchWeatherFromAddress(lat, lng float64) (*kakao.WeatherRequest, error) {
	return mfs.FacilityService.FetchWeatherFromAddress(lat, lng)
}

func (mfs *MarkerFacadeService) IsInSouthKoreaPrecisely(lat, lng float64) bool {
	return mfs.MapUtil.IsInSouthKoreaPrecisely(lat, lng)
}

func (mfs *MarkerFacadeService) SaveOfflineMap(lat, lng float64) (string, error) {
	return mfs.LocationService.SaveOfflineMap(lat, lng)
}

func (mfs *MarkerFacadeService) SaveOfflineMap2(lat, lng float64) (string, string, error) {
	return mfs.LocationService.SaveOfflineMap2(lat, lng)
}

func (mfs *MarkerFacadeService) TestDynamic(lat, lng, scale float64, width, height int64) {
	mfs.LocationService.TestDynamic(lat, lng, scale, width, height)
}
