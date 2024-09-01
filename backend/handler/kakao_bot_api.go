package handler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Alfex4936/chulbong-kr/facade"
	"github.com/Alfex4936/chulbong-kr/middleware"
	"github.com/Alfex4936/chulbong-kr/service"
	"github.com/gofiber/fiber/v2"

	k "github.com/Alfex4936/kakao"
)

type KakaoBotHandler struct {
	MarkerFacadeService *facade.MarkerFacadeService
	BleveSearchService  *service.BleveSearchService

	recentMarkersKey string
	searchMarkersKey string
}

// NewKakaoBotHandler creates a new KakaoBotHandler with dependencies injected
func NewKakaoBotHandler(facade *facade.MarkerFacadeService, bleve *service.BleveSearchService) *KakaoBotHandler {

	return &KakaoBotHandler{
		MarkerFacadeService: facade,
		BleveSearchService:  bleve,
		recentMarkersKey:    "kakaobot:recent-markers",
		searchMarkersKey:    "kakaobot:search-markers:",
	}
}

// RegisterKakaoBotRoutes sets up the routes for kakaotalk chatbot handling within the application.
func RegisterKakaoBotRoutes(api fiber.Router, handler *KakaoBotHandler, authMiddleware *middleware.AuthMiddleware) {
	kakaoGroup := api.Group("/kakaobot")
	{
		kakaoGroup.Post("/markers/recent-with-photos", handler.HandleKakaoRecentWithPhotos)
		kakaoGroup.Post("/markers/search", handler.HandleKakaoSearchMarkers)
	}
}

func (h *KakaoBotHandler) HandleKakaoRecentWithPhotos(c *fiber.Ctx) error {
	// Attempt to retrieve from cache first
	var response k.K
	cacheErr := h.MarkerFacadeService.GetRedisCache(h.recentMarkersKey, &response)
	if cacheErr == nil && response != nil {
		// Cache hit, return cached response
		return c.JSON(response)
	}

	// Fetch the top 10 markers with extra information
	markers, err := h.MarkerFacadeService.GetNew10PicturesWithExtra()
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(k.SimpleText{}.Build("잠시 후 다시 시도해주세요.", nil))
	}

	// Initialize the carousel without CommerceCard and CarouselHeader
	carousel := k.Carousel{}.New(false, false)

	for _, marker := range markers {
		// Create a new BasicCard with title, description, and thumbnail
		card := k.BasicCard{}.New(true, true)

		addresses := strings.Split(marker.Address, ",")

		card.Title = fmt.Sprintf("%s - %d", addresses[0], marker.MarkerID)

		address := strings.Split(marker.Address, ",")[1]

		var descBuilder strings.Builder
		descBuilder.WriteString("날씨: ")
		descBuilder.WriteString(marker.Weather)
		descBuilder.WriteString("\n" + address)
		card.Desc = descBuilder.String()

		thumbnail := k.ThumbNail{}.New(marker.PhotoURL)
		thumbnail.Link = &k.Link{Link: ("https://k-pullup.com/pullup/" + strconv.Itoa(marker.MarkerID))}
		card.ThumbNail = thumbnail

		// Add a link button to the card
		card.Buttons.Add(k.LinkButton{}.New("철봉 웹사이트", "https://k-pullup.com/pullup/"+strconv.Itoa(marker.MarkerID)))

		// Add the card to the carousel
		carousel.Cards.Add(card)
	}

	json := carousel.Build()

	go h.MarkerFacadeService.SetRedisCache(h.recentMarkersKey, json, 1*time.Hour)

	return c.JSON(json)
}

func (h *KakaoBotHandler) HandleKakaoSearchMarkers(c *fiber.Ctx) error {
	var kakaoRequest k.Request

	// Bind the request body
	if err := c.BodyParser(&kakaoRequest); err != nil {
		return c.Status(fiber.StatusOK).JSON(k.SimpleText{}.Build("잠시 후 다시 시도해주세요.", k.Kakao{
			k.QuickReply{}.New("철봉 검색", "검색"),
		}))
	}

	// Fetch
	utterance := kakaoRequest.Action.Params["search"].(string)
	utterance = strings.TrimSpace(utterance)

	// Attempt to retrieve from cache first
	var cacheResponse k.K
	cacheErr := h.MarkerFacadeService.GetRedisCache(h.searchMarkersKey+utterance, &cacheResponse)
	if cacheErr == nil && cacheResponse != nil {
		// Cache hit, return cached response
		return c.JSON(cacheResponse)
	}

	response, err := h.BleveSearchService.SearchMarkerAddress(utterance)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(k.SimpleText{}.Build("잠시 후 다시 시도해주세요.",
			k.Kakao{
				k.QuickReply{}.New("철봉 검색", "검색"),
			},
		))
	}

	if len(response.Markers) > 5 {
		response.Markers = response.Markers[:5]
	} else if len(response.Markers) == 0 {
		return c.Status(fiber.StatusOK).JSON(k.SimpleText{}.Build(utterance+"에 대한 검색 결과가 없습니다.", k.Kakao{
			k.QuickReply{}.New("철봉 검색", "검색"),
		}))
	}

	listCard := k.ListCard{}.New(true) // whether to use quickReplies or not
	listCard.Title = utterance + " 결과"

	for _, pullup := range response.Markers {
		listCard.Items.Add(k.ListItemLink{}.New(
			strconv.Itoa(pullup.MarkerID), pullup.Address, "",
			"https://k-pullup.com/pullup/"+strconv.Itoa(pullup.MarkerID)))
	}

	listCard.Buttons.Add(k.ShareButton{}.New("공유하기"))
	listCard.Buttons.Add(k.LinkButton{}.New("대한민국 철봉", "https://k-pullup.com"))

	// QuickReplies: label, message (메시지 없으면 라벨로 발화)
	listCard.QuickReplies.Add(k.QuickReply{}.New("철봉 검색", "검색"))
	listCard.QuickReplies.Add(k.QuickReply{}.New("최근 철봉", "최근"))

	cacheResponse = listCard.Build()

	go h.MarkerFacadeService.SetRedisCache(h.searchMarkersKey+utterance, cacheResponse, 1*time.Hour)

	return c.JSON(cacheResponse)
}
