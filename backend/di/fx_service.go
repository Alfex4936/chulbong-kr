package servicefx

import (
	"github.com/Alfex4936/chulbong-kr/facade"
	"github.com/Alfex4936/chulbong-kr/handler"
	"github.com/Alfex4936/chulbong-kr/service"
	"github.com/Alfex4936/chulbong-kr/util"
	"go.uber.org/fx"
)

var (
	FxAPIModule = fx.Module("api",
		fx.Provide(
			handler.NewMarkerHandler,
			handler.NewUserHandler,
			handler.NewSearchHandler,
			handler.NewNotificationHandler,
			handler.NewCommentHandler,
			handler.NewChatHandler,
			handler.NewAuthHandler,
			handler.NewAdminHandler,
		),
	)

	FxFacadeModule = fx.Module("facade",
		fx.Provide(
			facade.NewMarkerFacadeService,
			facade.NewUserFacadeService,
			facade.NewAdminFacadeService,
		),
	)

	FxUserModule = fx.Module("user",
		fx.Provide(
			service.NewAuthService,
			service.NewUserService,
		),
	)
	FxMarkerModule = fx.Module("marker",
		fx.Provide(
			service.NewMarkerManageService,
			service.NewMarkerInteractService,
			service.NewMarkerRankService,
			service.NewMarkerLocationService,
			service.NewMarkerFacilityService,
			service.NewMarkerCommentService,
			service.NewNotificationService,
			service.NewReportService,
		),
	)

	FxExternalModle = fx.Module("external",
		fx.Provide(
			service.NewRedisService,
			service.NewS3Service,
			service.NewZincSearchService,
			service.NewBleveSearchService,
			service.NewSmtpService,
		),
	)

	FxChatModule = fx.Module("chat",
		fx.Provide(
			service.NewChatService,
			service.NewRoomConnectionManager,
			// service.NewMqService,
			// service.NewGeminiService,
		),
	)

	FxUtilModule = fx.Module("util",
		fx.Provide(
			service.NewTokenService,
			service.NewSchedulerService,
			util.NewTokenUtil,
			util.NewChatUtil,
			util.NewBadWordUtil,
			util.NewMapUtil,
		),
	)
)
