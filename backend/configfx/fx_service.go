package configfx

import (
	"github.com/Alfex4936/chulbong-kr/config"
	"go.uber.org/fx"
)

var (
	FxConfigModule = fx.Module("config",
		fx.Provide(
			config.NewAppConfig,
			config.NewKakaoConfig,
			config.NewRedisConfig,
			config.NewZincSearchConfig,
			config.NewS3Config,
			config.NewSmtpConfig,
			config.NewTossPayConfig,
		),
	)
)
