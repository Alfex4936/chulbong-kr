package dto

type OAuthGoogleUser struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	ProfileImage string `json:"picture,omitempty"`
}

// KakaoProfile represents the profile information in the kakao_account object
type KakaoProfile struct {
	Nickname       string `json:"nickname"`
	ThumbnailImage string `json:"thumbnail_image_url,omitempty"`
	ProfileImage   string `json:"profile_image_url,omitempty"`
}

// KakaoAccount represents the kakao_account object
type KakaoAccount struct {
	Email   string       `json:"email"`
	Profile KakaoProfile `json:"profile"`
}

// OAuthKakaoUser represents the user information returned by Kakao
type OAuthKakaoUser struct {
	ID           int64        `json:"id"`
	KakaoAccount KakaoAccount `json:"kakao_account"`
	Properties   struct {
		Nickname string `json:"nickname"`
	} `json:"properties"`
}

type OAuthNaverUser struct {
	Response struct {
		ID           string `json:"id"`
		Email        string `json:"email"`
		Nickname     string `json:"nickname"`
		ProfileImage string `json:"profile_image,omitempty"`
	} `json:"response"`
}

type OAuthGitHubUser struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
