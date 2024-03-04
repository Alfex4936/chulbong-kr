package handlers

import (
	"encoding/base64"
	"os"

	"github.com/gofiber/fiber/v2"
)

var (
	TOSS_SECRET_KEY      = "Basic " + base64.StdEncoding.EncodeToString([]byte(os.Getenv("TOSS_SECRET_KEY_TEST")+":"))
	TOSS_CONFIRM_API_URL = "https://api.tosspayments.com/v1/payments/confirm"
	TOSS_PAYMENT_API_URL = "https://api.tosspayments.com/v1/payments/"
)

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

func ConfirmToss(c *fiber.Ctx) error {
	// Extract paymentKey, orderId, and amount from the request body
	var requestBody struct {
		PaymentKey string `json:"paymentKey"`
		OrderId    string `json:"orderId"`
		Amount     int    `json:"amount"`
	}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Prepare the Authorization header
	agent := fiber.Post(TOSS_CONFIRM_API_URL)
	agent.Set("Authorization", TOSS_SECRET_KEY)
	agent.ContentType("application/json")
	agent.Body(c.Body()) // json: {orderId: "aaa", amount: ?, paymentKey: "bbb"}

	// Send the request to the external API
	statusCode, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	// Forward the status code and body received from the external API
	return c.Status(statusCode).Send(body)
}
