package facade

import (
	"mime/multipart"
	"time"

	"github.com/Alfex4936/chulbong-kr/dto"
	"github.com/Alfex4936/chulbong-kr/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

// AdminFacadeService provides a simplified interface to various admin-related services.
type AdminFacadeService struct {
	MarkerManage   *service.MarkerManageService
	S3Service      *service.S3Service
	ChatService    *service.ChatService
	MarkerFacility *service.MarkerFacilityService
}

type AdminFacadeParams struct {
	fx.In

	MarkerManage   *service.MarkerManageService
	S3Service      *service.S3Service
	ChatService    *service.ChatService
	MarkerFacility *service.MarkerFacilityService
}

func NewAdminFacadeService(
	p AdminFacadeParams,
) *AdminFacadeService {
	return &AdminFacadeService{
		MarkerManage:   p.MarkerManage,
		S3Service:      p.S3Service,
		ChatService:    p.ChatService,
		MarkerFacility: p.MarkerFacility,
	}
}

func (mfs *AdminFacadeService) FetchAllPhotoURLsFromDB() ([]string, error) {
	return mfs.MarkerManage.FetchAllPhotoURLsFromDB()
}

func (mfs *AdminFacadeService) ListAllObjectsInS3() ([]string, error) {
	return mfs.S3Service.ListAllObjectsInS3()
}

func (mfs *AdminFacadeService) FindUnreferencedS3Objects(dbURLs []string, s3Keys []string) []string {
	return mfs.S3Service.FindUnreferencedS3Objects(dbURLs, s3Keys)
}

func (mfs *AdminFacadeService) DeleteDataFromS3(dataURL string) error {
	return mfs.S3Service.DeleteDataFromS3(dataURL)
}

func (mfs *AdminFacadeService) BanUser(markerID, userID string, duration time.Duration) error {
	return mfs.ChatService.BanUser(markerID, userID, duration)
}

func (mfs *AdminFacadeService) FetchLatestMarkers(thresholdDate time.Time) ([]service.DataItem, error) {
	return mfs.MarkerFacility.FetchLatestMarkers(thresholdDate)
}

func (mfs *AdminFacadeService) CheckMarkerValidity(latitude, longitude float64, description string) *fiber.Error {
	return mfs.MarkerManage.CheckMarkerValidity(latitude, longitude, description)
}

func (mfs *AdminFacadeService) CreateMarkerWithPhotos(markerDto *dto.MarkerRequest, userID int, form *multipart.Form) (*dto.MarkerResponse, error) {
	return mfs.MarkerManage.CreateMarkerWithPhotos(markerDto, userID, form)
}

func (mfs *AdminFacadeService) SetMarkerFacilities(markerID int, facilities []dto.FacilityQuantity) error {
	return mfs.MarkerFacility.SetMarkerFacilities(markerID, facilities)
}
