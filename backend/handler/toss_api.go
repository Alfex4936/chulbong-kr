package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Alfex4936/chulbong-kr/config"
	"github.com/gofiber/fiber/v2"
)

type TossHandler struct {
	PayConfig  *config.TossPayConfig
	HTTPClient *http.Client
}

// NewTossHandler creates a new TossHandler with dependencies injected
func NewTossHandler(payConfig *config.TossPayConfig, client *http.Client,
) *TossHandler {
	return &TossHandler{
		PayConfig:  payConfig,
		HTTPClient: client,
	}
}

// RegisterTossPaymentRoutes sets up the routes for toss payments handling within the application.
func RegisterTossPaymentRoutes(api fiber.Router, handler *TossHandler) {
	tossGroup := api.Group("/payments/toss")
	{
		tossGroup.Post("/confirm", handler.HandleConfirmToss)
		// tossGroup.Get("/success", handlers.SuccessToss)
		// tossGroup.Get("/fail", handlers.FailToss)
	}
}

// frontend handles
// func SuccessToss(c *fiber.Ctx) error {
// 	// Extract query parameters
// 	// orderId := c.Query("orderId")
// 	// amount := c.Query("amount")
// 	paymentKey := c.Query("paymentKey")

// 	// Prepare the request body
// 	// requestBody, _ := json.Marshal(map[string]interface{}{
// 	// 	"orderId":    orderId,
// 	// 	"amount":     amount,
// 	// 	"paymentKey": paymentKey,
// 	// })

// 	// Prepare the Authorization header
// 	agent := fiber.Get(TOSS_PAYMENT_API_URL + paymentKey)
// 	agent.Set("Authorization", TOSS_SECRET_KEY)
// 	agent.ContentType("application/json")
// 	// agent.Body(requestBody)

// 	// Send the request to the external API
// 	statusCode, body, errs := agent.Bytes()
// 	if len(errs) > 0 {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"errs": errs,
// 		})
// 	}

// 	var responseMap dto.Payment
// 	err := json.Unmarshal(body, &responseMap)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}

// 	log.Printf("Code: %+v", statusCode)
// 	log.Printf("errs: %+v", errs)
// 	log.Printf("Response: %+v", responseMap)

// 	// Render the template with response data
// 	return c.Render("success", fiber.Map{
// 		"IsSuccess": statusCode == 200,
// 		"Response":  responseMap,
// 	})
// }

// frontend handles
// func FailToss(c *fiber.Ctx) error {
// 	return c.SendString("Failed")
// }

func (h *TossHandler) HandleConfirmToss(c *fiber.Ctx) error {

	req, err := http.NewRequest(http.MethodPost, h.PayConfig.ConfirmAPI, bytes.NewBuffer(c.Body()))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("Authorization", h.PayConfig.SecretKey)
	resp, err := h.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var apiResp fiber.Map
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("unmarshalling response: %w", err)
	}

	// Forward the status code and body received from the external API
	return c.Status(resp.StatusCode).JSON(apiResp)
}
