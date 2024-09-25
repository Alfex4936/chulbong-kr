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

func (mfs *MarkerFacadeService) ApproveReport(reportID, userID int) error {
	if err := mfs.ReportService.ApproveReport(reportID, userID); err != nil {
		return err
	}

	return nil
}

func (mfs *MarkerFacadeService) DenyReport(reportID, userID int) error {
	return mfs.ReportService.DenyReport(reportID, userID)
}

func (mfs *MarkerFacadeService) DeleteReport(reportID, userID, markerID int) error {
	return mfs.ReportService.DeleteReport(reportID, userID, markerID)
}
