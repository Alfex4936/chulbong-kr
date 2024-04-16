package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
)

var Gemini *genai.Client

func ChatGemini() string {
	model := Gemini.GenerativeModel("gemini-1.0-pro")
	ctx := context.Background()
	cs := model.StartChat()

	send := func(msg string) *genai.GenerateContentResponse {
		res, _ := cs.SendMessage(ctx, genai.Text(msg))
		return res
	}

	res := send("Hello, how are you?")
	return formatResponse(res)
}

func formatResponse(resp *genai.GenerateContentResponse) string {
	var responseParts []string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				response := fmt.Sprintln(part)
				responseParts = append(responseParts, response)
			}
		}
	}
	return strings.Join(responseParts, "")
}
