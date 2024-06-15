package service

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Alfex4936/chulbong-kr/dto"
	"github.com/Alfex4936/chulbong-kr/model"
	"github.com/Alfex4936/chulbong-kr/protos"
	"github.com/Alfex4936/chulbong-kr/util"
	"github.com/gofiber/fiber/v2"

	gocache "github.com/eko/gocache/lib/v4/cache"

	ristretto_store "github.com/eko/gocache/store/ristretto/v4"

	"go.uber.org/fx"

	"github.com/jmoiron/sqlx"
)

// MarkerManageService is a service for marker management operations.
type MarkerManageService struct {
	DB *sqlx.DB

	MarkerLocationService *MarkerLocationService
	S3Service             *S3Service
	ZincSearchService     *ZincSearchService
	BleveSearchService    *BleveSearchService
	RedisService          *RedisService

	MapUtil           *util.MapUtil
	BadWordUtil       *util.BadWordUtil
	LocalCacheStorage *ristretto_store.RistrettoStore

	// markersLocalCache cache to store encoded marker data
	// markersLocalCache []byte // 400 kb is fine here
	// cacheMutex        sync.RWMutex

	byteCache *gocache.Cache[[]byte]
}

type MarkerManageServiceParams struct {
	fx.In

	DB                    *sqlx.DB
	MarkerLocationService *MarkerLocationService
	S3Service             *S3Service
	ZincSearchService     *ZincSearchService
	BleveSearchService    *BleveSearchService
	RedisService          *RedisService
	MapUtil               *util.MapUtil
	BadWordUtil           *util.BadWordUtil
	LocalCacheStorage     *ristretto_store.RistrettoStore
}

// NewMarkerManageService creates a new instance of MarkerManageService.
func NewMarkerManageService(p MarkerManageServiceParams) *MarkerManageService {
	byteCache := gocache.New[[]byte](p.LocalCacheStorage)

	return &MarkerManageService{
		DB:                    p.DB,
		MarkerLocationService: p.MarkerLocationService,
		S3Service:             p.S3Service,
		ZincSearchService:     p.ZincSearchService,
		BleveSearchService:    p.BleveSearchService,
		RedisService:          p.RedisService,
		MapUtil:               p.MapUtil,
		BadWordUtil:           p.BadWordUtil,
		byteCache:             byteCache,
	}
}

func RegisterMarkerLifecycle(lifecycle fx.Lifecycle, service *MarkerManageService) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			rss, _ := service.GenerateRSS()
			saveRSSToFile(rss, "marker_rss.xml")
			return nil
		},
		OnStop: func(context.Context) error {
			return nil
		},
	})
}

func (s *MarkerManageService) GetCache() []byte {
	// ctx := context.Background()
	// value, _ := s.byteCache.Get(ctx, "allMarkers")

	var value []byte
	s.RedisService.GetCacheEntry(s.RedisService.RedisConfig.AllMarkersKey, &value)
	return value
}

func (s *MarkerManageService) SetCache(mjson []byte) {
	s.RedisService.SetCacheEntry(s.RedisService.RedisConfig.AllMarkersKey, mjson, time.Hour*24)
	// ctx := context.Background()
	// s.byteCache.Set(ctx, "allMarkers", mjson)
}

func (s *MarkerManageService) ClearCache() {
	s.RedisService.ResetCache(s.RedisService.RedisConfig.AllMarkersKey)
	// ctx := context.Background()
	// s.byteCache.Delete(ctx, "allMarkers")
}

// GetAllMarkers now returns a simplified list of markers
func (s *MarkerManageService) GetAllMarkers() ([]dto.MarkerSimple, error) {
	// Simplified query to fetch only the marker IDs, latitudes, and longitudes
	const markerQuery = `
    SELECT 
        MarkerID, 
        ST_X(Location) AS Latitude,
        ST_Y(Location) AS Longitude
    FROM 
        Markers;`

	var markers []dto.MarkerSimple
	err := s.DB.Select(&markers, markerQuery)
	if err != nil {
		return nil, fmt.Errorf("error fetching markers: %w", err)
	}

	// go s.MarkerLocationService.Redis.AddGeoMarkers(markers)

	return markers, nil
}

// GetAllNewMarkers returns a paginated, simplified list of the most recently added markers.
func (s *MarkerManageService) GetAllNewMarkers(page, pageSize int) ([]dto.MarkerSimple, error) {
	if page < 1 {
		page = 1 // Ensure page starts at 1
	}
	if pageSize < 1 {
		pageSize = 10 // Default pageSize to 10 if invalid
	}
	offset := (page - 1) * pageSize

	const markerQuery = `
    SELECT 
        MarkerID, 
        ST_X(Location) AS Latitude,
        ST_Y(Location) AS Longitude
    FROM 
        Markers
    ORDER BY 
        createdAt DESC
    LIMIT ? OFFSET ?;`

	var markers []dto.MarkerSimple
	err := s.DB.Select(&markers, markerQuery, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("error fetching markers: %w", err)
	}

	return markers, nil
}

// GetMarker retrieves a single marker and its associated photo by the marker's ID
func (s *MarkerManageService) GetMarker(markerID int) (*model.MarkerWithPhotos, error) {
	const query = `
	SELECT
	M.MarkerID,
	M.UserID,
	ST_X(M.Location) AS Latitude,
	ST_Y(M.Location) AS Longitude,
	M.Description,
	COALESCE(U.Username, '탈퇴한 사용자') AS Username,
	M.CreatedAt,
	M.UpdatedAt,
	M.Address,
	COALESCE(D.DislikeCount, 0) AS DislikeCount,
	COALESCE(F.FavoriteCount, 0) AS FavoriteCount
  FROM Markers M
  LEFT JOIN Users U ON M.UserID = U.UserID
  LEFT JOIN (
	SELECT
	  MarkerID,
	  COUNT(*) AS DislikeCount
	FROM MarkerDislikes
	WHERE MarkerID = ?
	GROUP BY MarkerID
  ) D ON M.MarkerID = D.MarkerID
  LEFT JOIN (
	SELECT
	  MarkerID,
	  COUNT(*) AS FavoriteCount
	FROM Favorites
	WHERE MarkerID = ?
	GROUP BY MarkerID
  ) F ON M.MarkerID = F.MarkerID
  WHERE M.MarkerID = ?`

	var markersWithUsernames struct {
		model.Marker
		Username      string `db:"Username"`
		DislikeCount  int    `db:"DislikeCount"`
		FavoriteCount int    `db:"FavoriteCount"`
	}
	err := s.DB.Get(&markersWithUsernames, query, markerID, markerID, markerID)
	if err != nil {
		return nil, err
	}

	// Fetch all photos for this marker by descending order of upload time
	const photoQuery = `SELECT * FROM Photos WHERE MarkerID = ? ORDER BY UploadedAt DESC`
	var photos []model.Photo
	err = s.DB.Select(&photos, photoQuery, markerID)
	if err != nil {
		return nil, fmt.Errorf("error fetching photos: %w", err)
	}

	// Assemble the final structure
	markersWithPhotos := model.MarkerWithPhotos{
		Marker:        markersWithUsernames.Marker,
		Photos:        photos,
		Username:      markersWithUsernames.Username,
		DislikeCount:  markersWithUsernames.DislikeCount,
		FavoriteCount: markersWithUsernames.FavoriteCount,
	}

	// PublishMarkerUpdate(fmt.Sprintf("user: %s", markersWithPhotos.Username))

	return &markersWithPhotos, nil
}

func (s *MarkerManageService) GetAllMarkersProto() ([]*protos.Marker, error) {
	// Simplified query to fetch only the marker IDs, latitudes, and longitudes
	const markerQuery = `
    SELECT 
        MarkerID, 
        ST_X(Location) AS Latitude,
        ST_Y(Location) AS Longitude
    FROM 
        Markers;`

	var markers []*protos.Marker
	err := s.DB.Select(&markers, markerQuery)
	if err != nil {
		return nil, fmt.Errorf("error fetching markers: %w", err)
	}

	return markers, nil
}

// GetAllMarkersWithAddr fetches all markers and returns only those with an address not found or empty.
func (s *MarkerManageService) GetAllMarkersWithAddr() ([]dto.MarkerSimpleWithAddr, error) {
	const markerQuery = `SELECT MarkerID, ST_X(Location) AS Latitude, ST_Y(Location) AS Longitude FROM Markers;`

	var markers []dto.MarkerSimpleWithAddr
	err := s.DB.Select(&markers, markerQuery)
	if err != nil {
		return nil, fmt.Errorf("error fetching markers: %w", err)
	}

	var filteredMarkers []dto.MarkerSimpleWithAddr
	for i := range markers {
		// func FetchAddressFromAPI(latitude float64, longitude float64) (string, error)
		address, err := s.MarkerLocationService.FacilityService.FetchAddressFromAPI(markers[i].Latitude, markers[i].Longitude)
		if err != nil {
			continue
		}

		// Update the address and include only if "Address not found" or empty
		if address == "Address not found" || address == "" {
			markers[i].Address = address
			filteredMarkers = append(filteredMarkers, markers[i])
		}
	}

	return filteredMarkers, nil
}

func (s *MarkerManageService) GetAllMarkersByUserWithPagination(userID, page, pageSize int) ([]dto.MarkerSimpleWithDescrption, int, error) {
	offset := (page - 1) * pageSize

	// Query to select markers created by a specific user with LIMIT and OFFSET for pagination
	markerQuery := `
SELECT 
    M.MarkerID,
    ST_X(M.Location) AS Latitude, 
    ST_Y(M.Location) AS Longitude, 
    M.Description,
    M.CreatedAt,
	M.Address
FROM 
    Markers M
WHERE 
    M.UserID = ?
ORDER BY 
    M.CreatedAt DESC
LIMIT ? OFFSET ?
`

	markersWithDescription := make([]dto.MarkerSimpleWithDescrption, 0)
	err := s.DB.Select(&markersWithDescription, markerQuery, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	// Query to get the total count of markers for the user
	countQuery := `SELECT COUNT(DISTINCT Markers.MarkerID) FROM Markers WHERE Markers.UserID = ?`
	var total int
	err = s.DB.Get(&total, countQuery, userID)
	if err != nil {
		return nil, 0, err
	}

	return markersWithDescription, total, nil
}

func (s *MarkerManageService) CheckMarkerValidity(latitude, longitude float64, description string) *fiber.Error {
	var wg sync.WaitGroup
	errorChan := make(chan *fiber.Error, 3) // Buffered channel based on number of goroutines
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure all paths cancel the context to prevent context leak

	wg.Add(3) //three goroutines

	// South Korea check goroutine
	go func() {
		defer wg.Done()
		if inSKorea := s.MarkerLocationService.MapUtil.IsInSouthKoreaPrecisely(latitude, longitude); !inSKorea {
			select {
			case errorChan <- fiber.NewError(fiber.StatusForbidden, "operation is only allowed within South Korea"):
			case <-ctx.Done():
			}
			cancel() // Cancel other goroutines
		}
	}()

	// Marker proximity check goroutine
	go func() {
		defer wg.Done()
		if nearby, _ := s.MarkerLocationService.IsMarkerNearby(latitude, longitude, 10); nearby {
			select {
			case errorChan <- fiber.NewError(fiber.StatusConflict, "there is a marker already nearby"):
			case <-ctx.Done():
			}
			cancel() // Cancel other goroutines
		}
	}()

	// Bad words check goroutine
	go func() {
		defer wg.Done()
		if containsBadWord, _ := s.BadWordUtil.CheckForBadWordsUsingTrie(description); containsBadWord {
			select {
			case errorChan <- fiber.NewError(fiber.StatusBadRequest, "comment contains inappropriate content"):
			case <-ctx.Done():
			}
			cancel() // Cancel other goroutines
		}
	}()

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(errorChan) // Close channel safely after all sends are done
	}()

	// Process errors as they come in
	for err := range errorChan {
		return err
	}

	return nil
}

func (s *MarkerManageService) CreateMarkerWithPhotos(markerDto *dto.MarkerRequest, userID int, form *multipart.Form) (*dto.MarkerResponse, error) {
	s.ClearCache()

	// Begin a transaction for database operations
	tx, err := s.DB.Beginx()
	if err != nil {
		return nil, err
	}

	// Insert the marker into the database
	res, err := tx.Exec(
		"INSERT INTO Markers (UserID, Location, Description, CreatedAt, UpdatedAt) VALUES (?, ST_PointFromText(?, 4326), ?, NOW(), NOW())",
		userID, fmt.Sprintf("POINT(%f %f)", markerDto.Latitude, markerDto.Longitude), markerDto.Description,
	)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	markerID, _ := res.LastInsertId()
	folder := fmt.Sprintf("markers/%d", markerID)

	// After successfully creating the marker, process and upload the files
	var wg sync.WaitGroup
	errorChan := make(chan error, len(form.File["photos"]))

	// Process file uploads from the multipart form
	files := form.File["photos"]
	// Limit to a maximum of 5 files
	if len(files) > 5 {
		files = files[:5]
	}

	for _, file := range files {
		wg.Add(1)
		go func(file *multipart.FileHeader) {
			defer wg.Done()
			fileURL, err := s.S3Service.UploadFileToS3(folder, file)
			if err != nil {
				errorChan <- err
				return
			}
			if _, err := tx.Exec("INSERT INTO Photos (MarkerID, PhotoURL, UploadedAt) VALUES (?, ?, NOW())", markerID, fileURL); err != nil {
				errorChan <- err
				return
			}
		}(file)
	}

	// Wait for all uploads to complete
	wg.Wait()
	close(errorChan)

	// Check for errors
	for range errorChan {
		tx.Rollback()
		return nil, fmt.Errorf("encountered an error during file upload or DB operation")
	}

	// Commit the transaction after all operations succeed
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	go func(markerID int64, latitude, longitude float64) {
		maxRetries := 3
		retryDelay := 5 * time.Second // Delay between retries

		var address string
		var err error

		for attempt := 0; attempt < maxRetries; attempt++ {
			if attempt > 0 {
				log.Printf("Retrying to fetch address for marker %d, attempt %d", markerID, attempt)
				time.Sleep(retryDelay) // Wait before retrying
			}

			address, err = s.MarkerLocationService.FacilityService.FetchAddressFromMap(latitude, longitude)
			if err == nil && address != "" {
				break // Success, exit the retry loop
			}

			log.Printf("Attempt %d failed to fetch address for marker %d: %v", attempt, markerID, err)
		}

		if err != nil || address == "" {
			address2, _ := s.MarkerLocationService.FacilityService.FetchRegionFromAPI(latitude, longitude)
			if address2 == "-2" {
				// delete the marker 북한 or 일본
				_, err = s.DB.Exec("DELETE FROM Markers WHERE MarkerID = ?", markerID)
				if err != nil {
					log.Printf("Failed to update address for marker %d: %v", markerID, err)
				}
				return // no need to insert in failures
			}

			var errMsg string
			if err != nil {
				errMsg = fmt.Sprintf("Final attempt failed to fetch address for marker %d: %v", markerID, err)
			} else {
				errMsg = fmt.Sprintf("No address found for marker %d after %d attempts", markerID, maxRetries)
			}

			water, _ := s.MarkerLocationService.FacilityService.FetchRegionWaterInfo(latitude, longitude)
			if water {
				errMsg = fmt.Sprintf("The marker (%d) might be above on water", markerID)
			}

			url := fmt.Sprintf("%sd=%d&la=%f&lo=%f", s.MarkerLocationService.Config.ClientURL, markerID, markerDto.Latitude, markerDto.Longitude)

			logFailureStmt := "INSERT INTO MarkerAddressFailures (MarkerID, ErrorMessage, URL) VALUES (?, ?, ?)"
			if _, logErr := s.DB.Exec(logFailureStmt, markerID, errMsg, url); logErr != nil {
				log.Printf("Failed to log address fetch failure for marker %d: %v", markerID, logErr)
			}
			return
		}

		// Update the marker's address in the database after successful fetch
		_, err = s.DB.Exec("UPDATE Markers SET Address = ? WHERE MarkerID = ?", address, markerID)
		if err != nil {
			log.Printf("Failed to update address for marker %d: %v", markerID, err)
		}

		go s.BleveSearchService.InsertMarkerIndex(dto.MarkerIndexData{MarkerID: int(markerID), Address: address})
		// userIDstr := strconv.Itoa(userID)
		// updateMsg := fmt.Sprintf("새로운 철봉이 [ %s ]에 등록되었습니다!", address)
		// metadata := notification.NotificationMarkerMetadata{
		// 	MarkerID:  markerID,
		// 	Latitude:  latitude,
		// 	Longitude: longitude,
		// 	Address:   address,
		// }

		// rawMetadata, _ := json.Marshal(metadata)
		// PostNotification(userIDstr, "NewMarker", "k-pullup!", updateMsg, rawMetadata)

		// TODO: update when frontend updates
		// if address != "" {
		// 	updateMsg := fmt.Sprintf("새로운 철봉이 등록되었습니다! %s", address)
		// 	PublishMarkerUpdate(updateMsg)
		// }

	}(markerID, markerDto.Latitude, markerDto.Longitude)

	go s.MarkerLocationService.Redis.ResetAllCache(fmt.Sprintf("userMarkers:%d:page:*", userID))

	// Construct and return the response
	return &dto.MarkerResponse{
		MarkerID:    int(markerID),
		Latitude:    markerDto.Latitude,
		Longitude:   markerDto.Longitude,
		Description: markerDto.Description,
		// PhotoURLs:   photoURLs,
	}, nil
}

func (s *MarkerManageService) FetchAllPhotoURLsFromDB() ([]string, error) {
	query := `
        SELECT PhotoURL FROM Photos WHERE PhotoURL IS NOT NULL
        UNION
        SELECT URL FROM MarkerAddressFailures WHERE URL IS NOT NULL
        UNION
        SELECT ReportImageURL FROM Reports WHERE ReportImageURL IS NOT NULL
    `
	// Execute the query
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return urls, nil
}

// UpdateMarker updates an existing marker's latitude, longitude, and description
func (s *MarkerManageService) UpdateMarker(marker *model.Marker) error {
	const query = `UPDATE Markers SET Latitude = ?, Longitude = ?, Description = ?, UpdatedAt = NOW() 
                   WHERE MarkerID = ?`
	_, err := s.DB.Exec(query, marker.Latitude, marker.Longitude, marker.Description, marker.MarkerID)
	return err
}

func (s *MarkerManageService) UpdateMarkerDescriptionOnly(markerID int, description string) error {
	const query = `UPDATE Markers SET Description = ?, UpdatedAt = NOW() 
                   WHERE MarkerID = ?`
	_, err := s.DB.Exec(query, description, markerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no marker found with markerID %d", markerID)
		}
		return fmt.Errorf("error updating a marker: %w", err)
	}

	return nil
}

func (s *MarkerManageService) DeleteMarker(userID, markerID int, userRole string) error {
	// Precheck user authorization and fetch photo URLs before transaction
	var ownerID int
	var photoURLs []string
	const checkOwnerQuery = `SELECT UserID FROM Markers WHERE MarkerID = ?`
	const selectPhotosQuery = `SELECT PhotoURL FROM Photos WHERE MarkerID = ?`

	err := s.DB.Get(&ownerID, checkOwnerQuery, markerID)
	if err != nil {
		return fmt.Errorf("checking marker ownership: %w", err)
	}

	if userRole != "admin" && ownerID != userID {
		return fmt.Errorf("user %d is not authorized to delete marker %d", userID, markerID)
	}

	err = s.DB.Select(&photoURLs, selectPhotosQuery, markerID)
	if err != nil {
		return fmt.Errorf("fetching photos: %w", err)
	}

	// Start a transaction
	tx, err := s.DB.Beginx()
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	// Delete photos from database
	const deletePhotosQuery = `DELETE FROM Photos WHERE MarkerID = ?`
	if _, err := tx.Exec(deletePhotosQuery, markerID); err != nil {
		tx.Rollback()
		return fmt.Errorf("deleting photos: %w", err)
	}

	// Delete the marker
	const deleteMarkerQuery = `DELETE FROM Markers WHERE MarkerID = ?`
	if _, err := tx.Exec(deleteMarkerQuery, markerID); err != nil {
		tx.Rollback()
		return fmt.Errorf("deleting marker: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	// Delete photos from S3 in a goroutine
	go func(photoURLs []string) {
		for _, photoURL := range photoURLs {
			if err := s.S3Service.DeleteDataFromS3(photoURL); err != nil {
				log.Printf("Failed to delete photo from S3: %s, error: %v", photoURL, err)
			}
		}
	}(photoURLs)

	s.ClearCache()
	go s.BleveSearchService.DeleteMarkerIndex(strconv.Itoa(markerID))

	return nil
}

func (s *MarkerManageService) UploadMarkerPhotoToS3(markerID int, files []*multipart.FileHeader) ([]string, error) {
	// Begin a transaction for database operations
	tx, err := s.DB.Beginx()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	folder := fmt.Sprintf("markers/%d", markerID)

	picUrls := make([]string, 0)
	// Process file uploads from the multipart form
	for _, file := range files {
		fileURL, err := s.S3Service.UploadFileToS3(folder, file)
		if err != nil {
			fmt.Printf("Failed to upload file to S3: %v\n", err)
			continue // Skip this file and continue with the next
		}
		picUrls = append(picUrls, fileURL)
		// Associate each photo with the marker in the database
		if _, err := tx.Exec("INSERT INTO Photos (MarkerID, PhotoURL, UploadedAt) VALUES (?, ?, NOW())", markerID, fileURL); err != nil {
			// Attempt to delete the uploaded file from S3
			if delErr := s.S3Service.DeleteDataFromS3(fileURL); delErr != nil {
				fmt.Printf("Also failed to delete the file from S3: %v\n", delErr)
			}
			return nil, err
		}
	}

	// Update Marker's UpdatedAt field
	if _, err := tx.Exec("UPDATE Markers SET UpdatedAt = NOW() WHERE MarkerID = ?", markerID); err != nil {
		return nil, err
	}

	// If no errors, commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return picUrls, nil
}

func (s *MarkerManageService) CheckNearbyMarkersInDB() ([]dto.MarkerGroup, error) {
	markers, err := s.GetAllMarkers()
	if err != nil {
		return nil, fmt.Errorf("error fetching markers: %w", err)
	}

	var markerGroups []dto.MarkerGroup

	for _, marker := range markers {
		point := fmt.Sprintf("POINT(%f %f)", marker.Latitude, marker.Longitude)

		var nearbyMarkers []dto.MarkerWithDistance
		err := s.DB.Select(&nearbyMarkers, `
			SELECT MarkerID, ST_X(Location) AS Latitude, ST_Y(Location) AS Longitude, Description, ST_Distance_Sphere(Location, ST_GeomFromText(?, 4326)) AS distance, Address
			FROM Markers
			WHERE ST_Distance_Sphere(Location, ST_GeomFromText(?, 4326)) <= ?
			ORDER BY distance ASC
		`, point, point, 10)
		if err != nil {
			return nil, fmt.Errorf("error checking for nearby markers: %w", err)
		}

		// 필터링하여 자기 자신을 제외한 마커만 포함시키기
		var filteredNearbyMarkers []dto.MarkerWithDistance
		for _, nMarker := range nearbyMarkers {
			if nMarker.MarkerID != marker.MarkerID {
				filteredNearbyMarkers = append(filteredNearbyMarkers, nMarker)
			}
		}

		// 주변에 다른 마커들이 있는 경우에만 결과에 추가
		if len(filteredNearbyMarkers) > 0 {
			markerGroup := dto.MarkerGroup{
				CentralMarker: marker,
				NearbyMarkers: filteredNearbyMarkers,
			}
			markerGroups = append(markerGroups, markerGroup)
		}
	}

	return markerGroups, nil
}

func (s *MarkerManageService) GenerateRSS() (string, error) {
	const markerQuery = `SELECT MarkerID, UpdatedAt, Address FROM Markers ORDER BY UpdatedAt DESC`

	var markers []dto.MarkerRSS
	err := s.DB.Select(&markers, markerQuery)
	if err != nil {
		return "", fmt.Errorf("error fetching markers: %w", err)
	}

	return generateRSS(markers)
}

func generateRSS(markers []dto.MarkerRSS) (string, error) {
	items := []dto.RssItem{}
	for _, marker := range markers {
		item := dto.RssItem{
			Title:       fmt.Sprintf("Marker %d", marker.MarkerID),
			Link:        fmt.Sprintf("https://www.k-pullup.com/pullup/%d", marker.MarkerID),
			Description: marker.Address,
			PubDate:     marker.UpdatedAt.Format(time.RFC1123Z),
		}
		items = append(items, item)
	}

	rss := dto.RSS{
		Version: "2.0",
		Channel: dto.RssChannel{
			Title:         "k-pullup Markers",
			Link:          "https://www.k-pullup.com",
			Description:   "Latest markers of public pull-up bars",
			LastBuildDate: time.Now().Format(time.RFC1123Z),
			Items:         items,
		},
	}

	output, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshalling RSS: %w", err)
	}
	return `<?xml version="1.0" encoding="UTF-8"?>` + "\n" + string(output), nil
}

func saveRSSToFile(content, filePath string) error {
	return os.WriteFile(filePath, []byte(content), 0644)
}
