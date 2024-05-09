package facade

import (
	"mime/multipart"

	"github.com/Alfex4936/chulbong-kr/dto"
)

// Get
func (mfs *MarkerFacadeService) GetAllReports() ([]dto.MarkerReportResponse, error) {
	return mfs.ReportService.GetAllReports()
}

func (mfs *MarkerFacadeService) GetAllReportsBy(markerID int) ([]dto.MarkerReportResponse, error) {
	return mfs.ReportService.GetAllReportsBy(markerID)
}

func (mfs *MarkerFacadeService) CreateReport(report *dto.MarkerReportRequest, form *multipart.Form) error {
	return mfs.ReportService.CreateReport(report, form)
}
