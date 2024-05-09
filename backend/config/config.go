package config

import (
	"encoding/base64"
	"os"
	"strconv"
	"time"
)

type AppConfig struct {
	AwsRegion           string
	S3BucketName        string
	LoginTokenCookie    string
	TokenExpirationTime time.Duration
	IsProduction        string
	IsWaterURL          string
	IsWaterKEY          string
	CkURL               string // to fetch new ones from external
	NaverEmailVerifyURL string
	ClientURL           string
}

func NewAppConfig() *AppConfig {
	expiration, _ := strconv.Atoi(os.Getenv("TOKEN_EXPIRATION_INTERVAL"))
	return &AppConfig{
		AwsRegion:           os.Getenv("AWS_REGION"),
		S3BucketName:        os.Getenv("AWS_BUCKET_NAME"),
		LoginTokenCookie:    os.Getenv("TOKEN_COOKIE"),
		TokenExpirationTime: time.Duration(expiration) * time.Hour,
		IsProduction:        os.Getenv("DEPLOYMENT"),
		IsWaterURL:          os.Getenv("IS_WATER_API"),
		IsWaterKEY:          os.Getenv("IS_WATER_API_KEY"),
		CkURL:               os.Getenv("CK_URL"),
		NaverEmailVerifyURL: os.Getenv("NAVER_EMAIL_VERIFY_URL"),
		ClientURL:           os.Getenv("CLIENT_ADDR"),
	}
}

type KakaoConfig struct {
	KakaoAK           string
	KakaoStaticMap    string
	KakaoGeoCode      string
	KakaoCoord2Addr   string
	KakaoCoord2Region string
	KakaoWeather      string
	KakaoAddressInfo  string
}

func NewKakaoConfig() *KakaoConfig {
	return &KakaoConfig{
		KakaoAK:           os.Getenv("KAKAO_AK"),
		KakaoStaticMap:    os.Getenv("KAKAO_STATIC_MAP"),
		KakaoGeoCode:      "https://dapi.kakao.com/v2/local/geo/address.json",
		KakaoCoord2Addr:   "https://dapi.kakao.com/v2/local/geo/coord2address.json",
		KakaoCoord2Region: "https://dapi.kakao.com/v2/local/geo/coord2regioncode.json",
		KakaoWeather:      "https://map.kakao.com/api/dapi/point/weather?inputCoordSystem=WCONGNAMUL&outputCoordSystem=WCONGNAMUL&version=2&service=map.daum.net",
		KakaoAddressInfo:  "https://map.kakao.com/etc/areaAddressInfo.json?output=JSON&inputCoordSystem=WCONGNAMUL&outputCoordSystem=WCONGNAMUL",
	}
}

type RedisConfig struct {
	AllMarkersKey  string
	UserProfileKey string
	UserFavKey     string
}

func NewRedisConfig() *RedisConfig {
	return &RedisConfig{
		AllMarkersKey:  "all_markers",
		UserProfileKey: "user_profile",
		UserFavKey:     "user_fav",
	}
}

type ZincSearchConfig struct {
	ZincAPI      string
	ZincUser     string
	ZincPassword string
}

func NewZincSearchConfig() *ZincSearchConfig {
	return &ZincSearchConfig{
		ZincAPI:      os.Getenv("ZINCSEARCH_URL"),
		ZincUser:     os.Getenv("ZINCSEARCH_USER"),
		ZincPassword: os.Getenv("ZINCSEARCH_PASSWORD"),
	}
}

type S3Config struct {
	AwsRegion    string
	S3BucketName string
}

func NewS3Config() *S3Config {
	return &S3Config{
		AwsRegion:    os.Getenv("AWS_REGION"),
		S3BucketName: os.Getenv("AWS_BUCKET_NAME"),
	}
}

type SmtpConfig struct {
	SmtpServer          string
	SmtpPort            string
	SmtpUsername        string
	SmtpPassword        string
	FrontendResetRouter string
}

func NewSmtpConfig() *SmtpConfig {
	return &SmtpConfig{
		SmtpServer:          os.Getenv("SMTP_SERVER"),
		SmtpPort:            os.Getenv("SMTP_PORT"),
		SmtpUsername:        os.Getenv("SMTP_USERNAME"),
		SmtpPassword:        os.Getenv("SMTP_PASSWORD"),
		FrontendResetRouter: os.Getenv("FRONTEND_RESET_ROUTER"),
	}
}

type TossPayConfig struct {
	SecretKey  string
	ConfirmAPI string
	PaymentAPI string
}

func NewTossPayConfig() *TossPayConfig {
	return &TossPayConfig{
		SecretKey:  "Basic " + base64.StdEncoding.EncodeToString([]byte(os.Getenv("TOSS_SECRET_KEY_TEST")+":")),
		ConfirmAPI: "https://api.tosspayments.com/v1/payments/confirm",
		PaymentAPI: "https://api.tosspayments.com/v1/payments/",
	}
}
