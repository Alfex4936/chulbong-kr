package facade

import (
	"github.com/Alfex4936/chulbong-kr/dto"
)

// Get
func (mfs *MarkerFacadeService) GetTopMarkers(limit int) []dto.MarkerSimpleWithAddr {
	return mfs.RankService.GetTopMarkers(limit)
}

func (mfs *MarkerFacadeService) GetUniqueVisitorCount(markerID string) int {
	return mfs.RankService.GetUniqueVisitorCount(markerID)
}

func (mfs *MarkerFacadeService) GetAllUniqueVisitorCounts() map[string]int {
	return mfs.RankService.GetAllUniqueVisitorCounts()
}
