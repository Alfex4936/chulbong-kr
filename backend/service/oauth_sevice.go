package service

import (
	"context"
	"io"
	"net/http"

	"github.com/Alfex4936/chulbong-kr/model"
	sonic "github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

func (s *AuthService) GoogleCallback(c *fiber.Ctx) (*model.User, error) {
	code := c.Query("code")
	token, err := s.OAuthConfig.GoogleOAuth.Exchange(context.TODO(), code)
	if err != nil {
		return nil, c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	client := s.OAuthConfig.GoogleOAuth.Client(context.TODO(), token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	userInfo := struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}{}
	if err := sonic.Unmarshal(body, &userInfo); err != nil {
		return nil, c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	user, err := s.SaveOAuthUser("google", userInfo.ID, userInfo.Email, userInfo.Name)
	if err != nil {
		return nil, c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return user, nil
}

func (s *AuthService) KakaoCallback(c *fiber.Ctx) (*model.User, error) {
	code := c.Query("code")
	token, err := s.OAuthConfig.KakaoOAuth.Exchange(context.TODO(), code)
	if err != nil {
		return nil, c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	client := s.OAuthConfig.KakaoOAuth.Client(context.TODO(), token)
	response, err := client.Get("https://kapi.kakao.com/v2/user/me")
	if err != nil {
		return nil, c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	userInfo := struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		Nickname string `json:"properties.nickname"`
	}{}
	if err := sonic.Unmarshal(body, &userInfo); err != nil {
		return nil, c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	user, err := s.SaveOAuthUser("kakao", userInfo.ID, userInfo.Email, userInfo.Nickname)
	if err != nil {
		return nil, c.Status(http.StatusInternalServerError).SendString(err.Error())

	}

	return user, nil
}
