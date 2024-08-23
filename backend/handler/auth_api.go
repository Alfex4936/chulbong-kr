package handler

import (
	"context"
	"database/sql"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/Alfex4936/chulbong-kr/dto"
	"github.com/Alfex4936/chulbong-kr/middleware"
	"github.com/Alfex4936/chulbong-kr/service"
	"github.com/Alfex4936/chulbong-kr/util"
	sonic "github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	GOOGLE_USER_INFO_URL = "https://www.googleapis.com/oauth2/v2/userinfo"
	KAKAO_USER_INFO_URL  = "https://kapi.kakao.com/v2/user/me"
	NAVER_USER_INFO_URL  = "https://openapi.naver.com/v1/nid/me"
)

type AuthHandler struct {
	AuthService  *service.AuthService
	UserService  *service.UserService
	TokenService *service.TokenService
	SmtpService  *service.SmtpService

	TokenUtil *util.TokenUtil
	Logger    *zap.Logger
}

// NewAuthHandler creates a new AuthHandler with dependencies injected
func NewAuthHandler(
	auth *service.AuthService,
	user *service.UserService,
	token *service.TokenService,
	smtp *service.SmtpService,
	tutil *util.TokenUtil,
	logger *zap.Logger,
) *AuthHandler {
	return &AuthHandler{
		AuthService:  auth,
		UserService:  user,
		TokenService: token,
		SmtpService:  smtp,
		TokenUtil:    tutil,
		Logger:       logger,
	}
}

// RegisterAuthRoutes sets up the routes for auth handling within the application.
func RegisterAuthRoutes(api fiber.Router, handler *AuthHandler, authMiddleaware *middleware.AuthMiddleware) {
	authGroup := api.Group("/auth")
	{
		// OAuth2
		authGroup.Get("/google", handler.HandleGoogleOAuth)
		authGroup.Get("/naver", handler.HandleNaverOAuth)
		authGroup.Get("/kakao", handler.HandleKakaoOAuth)
		authGroup.Get("/github", handler.HandleGitHubOAuth)

		authGroup.Post("/signup", handler.HandleSignUp)
		authGroup.Post("/login", handler.HandleLogin)
		authGroup.Post("/logout", authMiddleaware.Verify, handler.HandleLogout)
		authGroup.Post("/verify-email/send", handler.HandleSendVerificationEmail)
		authGroup.Post("/verify-email/confirm", handler.HandleValidateToken)

		// Finding password
		authGroup.Post("/request-password-reset", handler.HandleRequestResetPassword)
		authGroup.Post("/reset-password", handler.HandleResetPassword)
	}
}

// func (h *AuthHandler) HandleGoogleLogin(c *fiber.Ctx) error {
// 	url := h.AuthService.OAuthConfig.GoogleOAuth.AuthCodeURL("state")
// 	c.Status(fiber.StatusSeeOther)
// 	return c.Redirect(url)
// }

// func (h *AuthHandler) HandleKakaoLogin(c *fiber.Ctx) error {
// 	url := h.AuthService.OAuthConfig.KakaoOAuth.AuthCodeURL("state")
// 	c.Status(fiber.StatusSeeOther)
// 	return c.Redirect(url)
// }

// SignUp User godoc
//
//	@Summary		Sign up a new user [normal]
//	@Description	This endpoint is responsible for registering a new user in the system.
//	@Description	It checks the verification status of the user's email before proceeding.
//	@Description	If the email is not verified, it returns an error.
//	@Description	On successful creation, it returns the user's information.
//	@ID				sign-up-user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			signUpRequest	body	dto.SignUpRequest	true	"SignUp Request"
//	@Security		ApiKeyAuth
//	@Success		201	{object}	models.User		"User registered successfully"
//	@Failure		400	{object}	map[string]interface{}	"Cannot parse JSON, wrong sign up form."
//	@Failure		400	{object}	map[string]interface{}	"Email not verified"
//	@Failure		409	{object}	map[string]interface{}	"Email already registered"
//	@Failure		500	{object}	map[string]interface{}	"An error occurred while creating the user"
//	@Router			/auth/signup [post]
func (h *AuthHandler) HandleSignUp(c *fiber.Ctx) error {
	var signUpReq dto.SignUpRequest
	if err := c.BodyParser(&signUpReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON, wrong sign up form."})
	}

	// Check if the token is verified before proceeding
	verified, err := h.TokenService.IsTokenVerified(signUpReq.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check verification status"})
	}
	if !verified {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email not verified"})
	}

	signUpReq.Provider = "website"

	user, err := h.AuthService.SaveUser(&signUpReq)
	if err != nil {
		// Handle the duplicate email error
		if strings.Contains(err.Error(), "already registered") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		// For other errors, return a generic error message
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "An error occurred while creating the user: " + err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// Login User godoc
//
// @Summary		Log in a user
// @Description	This endpoint is responsible for authenticating a user in the system.
// @Description	It validates the user's login credentials (email and password).
// @Description	If the credentials are invalid, it returns an error.
// @Description	On successful authentication, it returns the user's information along with a token.
// @Description	The token is also set in a secure cookie for client-side storage.
// @ID			login-user
// @Tags		auth
// @Accept		json
// @Produce	json
// @Param		loginRequest	body	dto.LoginRequest	true	"Login Request"
// @Security	ApiKeyAuth
// @Success	200	{object}	dto.LoginResponse	"User logged in successfully, includes user info and token"
// @Failure	400	{object}	map[string]interface{}	"Cannot parse JSON, wrong login form."
// @Failure	401	{object}	map[string]interface{}	"Invalid email or password"
// @Failure	500	{object}	map[string]interface{}	"Failed to generate token"
// @Router		/auth/login [post]
func (h *AuthHandler) HandleLogin(c *fiber.Ctx) error {
	var request dto.LoginRequest
	if err := c.BodyParser(&request); err != nil {
		h.Logger.Error("Failed to parse login request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := h.AuthService.Login(request.Email, request.Password)
	if err != nil {
		h.Logger.Warn("Login failed", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	token, err := h.TokenService.GenerateAndSaveToken(user.UserID)
	if err != nil {
		h.Logger.Error("Failed to generate token", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	// Create a response object that includes both the user and the token
	response := dto.LoginResponse{
		User:  user,
		Token: token,
	}

	// Setting the token in a secure cookie
	cookie := h.TokenUtil.GenerateLoginCookie(token)
	c.Cookie(&cookie)

	h.Logger.Info("User logged in successfully", zap.Int("userID", user.UserID))

	return c.JSON(response)
}

func (h *AuthHandler) HandleGoogleOAuth(c *fiber.Ctx) error {
	// Check if the state and code are present in the query params
	state := c.Query("state")
	code := c.Query("code")

	if state == "" || code == "" {
		// Start the OAuth flow
		state, err := h.TokenUtil.GenerateOpaqueToken(16)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to generate state")
		}

		// Store the state value in the session or a temporary store
		c.Cookie(&fiber.Cookie{
			Name:  "oauth_state",
			Value: state,
		})

		url := h.AuthService.OAuthConfig.GoogleOAuth.AuthCodeURL(state)
		c.Status(fiber.StatusSeeOther)
		return c.Redirect(url)
	}

	// Validate the state parameter
	storedState := c.Cookies("oauth_state")
	if state != storedState {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid state parameter")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Exchange the code for a token
	otoken, err := h.AuthService.OAuthConfig.GoogleOAuth.Exchange(ctx, code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	client := h.AuthService.OAuthConfig.GoogleOAuth.Client(ctx, otoken)
	response, err := client.Get(GOOGLE_USER_INFO_URL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	var userInfo dto.OAuthGoogleUser
	if err := sonic.Unmarshal(body, &userInfo); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	user, err := h.AuthService.SaveOAuthUser("google", userInfo.ID, userInfo.Email, userInfo.Name)
	if err != nil {
		h.Logger.Warn("OAuth Login failed", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Failed to login with Google"})
	}

	token, err := h.TokenService.GenerateAndSaveToken(user.UserID)
	if err != nil {
		h.Logger.Error("Failed to generate token for google login", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	// responseData := dto.LoginResponse{
	// 	User:  user,
	// 	Token: token,
	// }

	cookie := h.TokenUtil.GenerateLoginCookie(token)
	c.Cookie(&cookie)

	h.Logger.Info("Google user logged in successfully", zap.Int("userID", user.UserID))

	customRedirectUrl := c.Query("returnUrl")
	if customRedirectUrl == "" {
		return c.Redirect(h.AuthService.OAuthConfig.FrontendURL + "/mypage")
	} else {
		return c.Redirect(h.AuthService.OAuthConfig.FrontendURL + customRedirectUrl) // should be like "/pullup/1234"
	}
}

func (h *AuthHandler) HandleKakaoOAuth(c *fiber.Ctx) error {
	// Check if the state and code are present in the query params
	state := c.Query("state")
	code := c.Query("code")

	if state == "" || code == "" {
		// Start the OAuth flow
		state, err := h.TokenUtil.GenerateOpaqueToken(16)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to generate state")
		}

		// Store the state value in the session or a temporary store
		c.Cookie(&fiber.Cookie{
			Name:  "oauth_state",
			Value: state,
		})

		url := h.AuthService.OAuthConfig.KakaoOAuth.AuthCodeURL(state)
		c.Status(fiber.StatusSeeOther)
		return c.Redirect(url)
	}

	// Validate the state parameter
	storedState := c.Cookies("oauth_state")
	if state != storedState {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid state parameter")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Exchange the code for a token
	otoken, err := h.AuthService.OAuthConfig.KakaoOAuth.Exchange(ctx, code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	client := h.AuthService.OAuthConfig.KakaoOAuth.Client(ctx, otoken)
	response, err := client.Get(KAKAO_USER_INFO_URL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	var userInfo dto.OAuthKakaoUser
	if err := sonic.Unmarshal(body, &userInfo); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	user, err := h.AuthService.SaveOAuthUser("kakao", strconv.FormatInt(userInfo.ID, 10), userInfo.KakaoAccount.Email, userInfo.KakaoAccount.Profile.Nickname)
	if err != nil {
		h.Logger.Warn("OAuth Kakao Login failed", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Failed to login with Google"})
	}

	token, err := h.TokenService.GenerateAndSaveToken(user.UserID)
	if err != nil {
		h.Logger.Error("Failed to generate token for kakao login", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	cookie := h.TokenUtil.GenerateLoginCookie(token)
	c.Cookie(&cookie)

	h.Logger.Info("Kakao user logged in successfully", zap.Int("userID", user.UserID))

	// Redirect to the frontend with a token
	customRedirectUrl := c.Query("returnUrl")
	if customRedirectUrl == "" {
		return c.Redirect(h.AuthService.OAuthConfig.FrontendURL + "/mypage")
	} else {
		return c.Redirect(h.AuthService.OAuthConfig.FrontendURL + customRedirectUrl) // should be like "/pullup/1234"
	}
}

func (h *AuthHandler) HandleNaverOAuth(c *fiber.Ctx) error {
	state := c.Query("state")
	code := c.Query("code")

	if state == "" || code == "" {
		state, err := h.TokenUtil.GenerateOpaqueToken(16)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to generate state")
		}

		c.Cookie(&fiber.Cookie{
			Name:  "oauth_state",
			Value: state,
		})

		url := h.AuthService.OAuthConfig.NaverOAuth.AuthCodeURL(state)
		c.Status(fiber.StatusSeeOther)
		return c.Redirect(url)
	}

	storedState := c.Cookies("oauth_state")
	if state != storedState {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid state parameter")
	}

	otoken, err := h.AuthService.OAuthConfig.NaverOAuth.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	client := h.AuthService.OAuthConfig.NaverOAuth.Client(context.Background(), otoken)
	response, err := client.Get(NAVER_USER_INFO_URL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	var userInfo dto.OAuthNaverUser
	if err := sonic.Unmarshal(body, &userInfo); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	user, err := h.AuthService.SaveOAuthUser("naver", userInfo.Response.ID, userInfo.Response.Email, userInfo.Response.Nickname)
	if err != nil {
		h.Logger.Warn("OAuth Login failed", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Failed to login with Naver"})
	}

	token, err := h.TokenService.GenerateAndSaveToken(user.UserID)
	if err != nil {
		h.Logger.Error("Failed to generate token for Naver login", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	cookie := h.TokenUtil.GenerateLoginCookie(token)
	c.Cookie(&cookie)

	h.Logger.Info("Naver user logged in successfully", zap.Int("userID", user.UserID))

	customRedirectUrl := c.Query("returnUrl")
	if customRedirectUrl == "" {
		return c.Redirect(h.AuthService.OAuthConfig.FrontendURL + "/mypage")
	} else {
		return c.Redirect(h.AuthService.OAuthConfig.FrontendURL + customRedirectUrl) // should be like "/pullup/1234"
	}
}

func (h *AuthHandler) HandleGitHubOAuth(c *fiber.Ctx) error {
	// Check if the state and code are present in the query params
	state := c.Query("state")
	code := c.Query("code")

	if state == "" || code == "" {
		// Start the OAuth flow
		state, err := h.TokenUtil.GenerateOpaqueToken(16)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to generate state")
		}

		// Store the state value in the session or a temporary store
		c.Cookie(&fiber.Cookie{
			Name:  "oauth_state",
			Value: state,
		})

		url := h.AuthService.OAuthConfig.GitHubOAuth.AuthCodeURL(state)
		c.Status(fiber.StatusSeeOther)
		return c.Redirect(url)
	}

	// Validate the state parameter
	storedState := c.Cookies("oauth_state")
	if state != storedState {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid state parameter")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Exchange the code for a token
	otoken, err := h.AuthService.OAuthConfig.GitHubOAuth.Exchange(ctx, code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	client := h.AuthService.OAuthConfig.GitHubOAuth.Client(ctx, otoken)

	// Fetch user information from GitHub API
	response, err := client.Get("https://api.github.com/user")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	var userInfo struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
		Email     string `json:"email"`
	}
	if err := sonic.Unmarshal(body, &userInfo); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// If Email is empty, we may need to fetch it separately because GitHub doesn't always include it in the main user profile response.
	if userInfo.Email == "" {
		emailResp, err := client.Get("https://api.github.com/user/emails")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		defer emailResp.Body.Close()

		emailBody, err := io.ReadAll(emailResp.Body)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		var emails []struct {
			Email    string `json:"email"`
			Primary  bool   `json:"primary"`
			Verified bool   `json:"verified"`
		}
		if err := sonic.Unmarshal(emailBody, &emails); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		// Look for the primary and verified email
		for _, email := range emails {
			if email.Primary && email.Verified {
				userInfo.Email = email.Email
				break
			}
		}
	}

	user, err := h.AuthService.SaveOAuthUser("github", strconv.FormatInt(userInfo.ID, 10), userInfo.Email, userInfo.Login)
	if err != nil {
		h.Logger.Warn("OAuth GitHub Login failed", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Failed to login with GitHub"})
	}

	token, err := h.TokenService.GenerateAndSaveToken(user.UserID)
	if err != nil {
		h.Logger.Error("Failed to generate token for GitHub login", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	cookie := h.TokenUtil.GenerateLoginCookie(token)
	c.Cookie(&cookie)

	h.Logger.Info("GitHub user logged in successfully", zap.Int("userID", user.UserID))

	// Redirect to the frontend with a token
	customRedirectUrl := c.Query("returnUrl")
	if customRedirectUrl == "" {
		return c.Redirect(h.AuthService.OAuthConfig.FrontendURL + "/mypage")
	} else {
		return c.Redirect(h.AuthService.OAuthConfig.FrontendURL + customRedirectUrl) // should be like "/pullup/1234"
	}
}

func (h *AuthHandler) HandleLogout(c *fiber.Ctx) error {
	// Retrieve user ID from context or session
	userID, ok := c.Locals("userID").(int)
	if !ok {
		// If user ID is missing, continue with logout without error
		// Log the occurrence for internal tracking
		h.Logger.Warn("UserID missing in session during logout")
	}

	token := c.Cookies("token")

	if token != "" {
		// Attempt to delete the token from the database
		if err := h.TokenService.DeleteOpaqueToken(userID, token); err != nil {
			// Log the error but do not disrupt the logout process
			h.Logger.Warn("Failed to delete session token", zap.Int("userID", userID), zap.Error(err))
		}
	}

	// Clear the authentication cookie
	cookie := h.TokenUtil.ClearLoginCookie()
	c.Cookie(&cookie)

	h.Logger.Info("User logged out successfully", zap.Int("userID", userID))

	// Return a logout success response regardless of server-side token deletion status
	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}

// Send Verification Email godoc
//
// @Summary		Send verification email
// @Description	This endpoint triggers sending a verification email to the user.
// @Description	It checks if the email is already registered in the system.
// @Description	If the email is already in use, it returns an error.
// @Description	If the email is not in use, it asynchronously sends a verification email to the user.
// @Description	The operation of sending the email does not block the API response, making use of a goroutine for asynchronous execution.
// @ID			send-verification-email
// @Tags		auth
// @Accept		json
// @Produce	json
// @Param		email	formData	string	true	"User Email"
// @Security	ApiKeyAuth
// @Success	200	"Email sending initiated successfully"
// @Failure	409	{object}	map[string]interface{}	"Email already registered"
// @Failure	500	{object}	map[string]interface{}	"An unexpected error occurred"
// @Router		/auth/send-verification-email [post]
func (h *AuthHandler) HandleSendVerificationEmail(c *fiber.Ctx) error {
	userEmail := c.FormValue("email")
	userEmail = strings.ToLower(userEmail)
	_, err := h.UserService.GetUserByEmail(userEmail)
	if err == nil {
		// If GetUserByEmail does not return an error, it means the email is already in use
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already registered"})
	} else if err != sql.ErrNoRows {
		// if db couldn't find a user, then it's valid. other errors are bad.
		h.Logger.Error("Unexpected error occurred while checking email", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "An unexpected error occurred"})
	}

	// No matter if it's verified, send again.
	// Check if there's already a verified token for this user
	// verified, err := services.IsTokenVerified(userEmail)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check verification status"})
	// }
	// if verified {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email already verified"})
	// }

	// token, err := services.GenerateAndSaveSignUpToken(userEmail)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	// }

	// Use a goroutine to send the email without blocking
	go func(email string) {
		if strings.HasSuffix(email, "@naver.com") { // endsWith
			exist, _ := h.AuthService.VerifyNaverEmail(email)
			if !exist {
				h.Logger.Warn("No such email found on Naver", zap.String("email", email))
				return
			}
		}
		token, err := h.TokenService.GenerateAndSaveSignUpToken(email)
		if err != nil {
			h.Logger.Error("Failed to generate token", zap.String("email", email), zap.Error(err))
			return
		}

		err = h.SmtpService.SendVerificationEmail(email, token)
		if err != nil {
			h.Logger.Error("Failed to send verification email", zap.String("email", email), zap.Error(err))
			return
		}
	}(userEmail)

	return c.SendStatus(fiber.StatusOK)
}

// Validate Token godoc
//
// @Summary		Validate token
// @Description	This endpoint is responsible for validating a user's token.
// @Description	It checks the token's validity against the provided email.
// @Description	If the token is invalid or expired, it returns an error.
// @Description	On successful validation, it returns a success status.
// @ID			validate-token
// @Tags		auth
// @Accept		json
// @Produce	json
// @Param		token	formData	string	true	"Token for validation"
// @Param		email	formData	string	true	"User's email associated with the token"
// @Security	ApiKeyAuth
// @Success	200	"Token validated successfully"
// @Failure	400	{object}	map[string]interface{}	"Invalid or expired token"
// @Failure	500	{object}	map[string]interface{}	"Error validating token"
// @Router		/auth/validate-token [post]
func (h *AuthHandler) HandleValidateToken(c *fiber.Ctx) error {
	token := c.FormValue("token")
	email := c.FormValue("email")

	valid, err := h.TokenService.ValidateToken(token, email)
	if err != nil {
		// If err is not nil, it could be a database error or token not found
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error validating token"})
	}
	if !valid {
		// Handle both not found and expired cases
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	return c.SendStatus(fiber.StatusOK)
}

// Request Reset Password godoc
//
// @Summary		Request password reset
// @Description	This endpoint initiates the password reset process for a user.
// @Description	It generates a password reset token and sends a reset password email to the user.
// @Description	The email sending process is executed in a non-blocking manner using a goroutine.
// @Description	If there is an issue generating the token or sending the email, it returns an error.
// @ID			request-reset-password
// @Tags		auth
// @Accept		json
// @Produce	json
// @Param		email	formData	string	true	"User's email address for password reset"
// @Security	ApiKeyAuth
// @Success	200	"Password reset request initiated successfully"
// @Failure	500	{object}	map[string]interface{}	"Failed to request reset password"
// @Router		/auth/request-reset-password [post]
func (h *AuthHandler) HandleRequestResetPassword(c *fiber.Ctx) error {
	email := c.FormValue("email")

	// Generate the password reset token
	token, err := h.AuthService.GeneratePasswordResetToken(email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to request reset password: " + err.Error()})
	}

	// Use a goroutine to send the email without blocking
	go func(email string) {

		// Send the reset email
		err = h.SmtpService.SendPasswordResetEmail(email, token)
		if err != nil {
			// cannot respond to the client at this point
			h.Logger.Error("Error sending reset email", zap.String("email", email), zap.Error(err))
			return
		}
	}(email)

	return c.SendStatus(fiber.StatusOK)
}

// Reset Password godoc
//
// @Summary		Reset password
// @Description	This endpoint allows a user to reset their password using a valid token.
// @Description	The token is typically obtained from a password reset email.
// @Description	If the token is invalid or the reset fails, it returns an error.
// @Description	On successful password reset, it returns a success status.
// @ID			reset-password
// @Tags		auth
// @Accept		json
// @Produce	json
// @Param		token		formData	string	true	"Password reset token"
// @Param		password	formData	string	true	"New password"
// @Security	ApiKeyAuth
// @Success	200	"Password reset successfully"
// @Failure	500	{object}	map[string]interface{}	"Failed to reset password"
// @Router		/auth/reset-password [post]
func (h *AuthHandler) HandleResetPassword(c *fiber.Ctx) error {
	token := c.FormValue("token")
	newPassword := c.FormValue("password")

	err := h.AuthService.ResetPassword(token, newPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to reset password"})
	}

	return c.SendStatus(fiber.StatusOK)
}
