package services

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

func CronCleanUpToken() {
	c := cron.New()
	_, err := c.AddFunc("@daily", func() {
		if err := DeleteExpiredTokens(); err != nil {
			// Log the error
			fmt.Printf("Error deleting expired tokens: %v\n", err)
		} else {
			fmt.Println("Expired tokens cleanup executed successfully")
		}
	})
	if err != nil {
		// Handle the error
		fmt.Printf("Error scheduling the token cleanup job: %v\n", err)
		return
	}
	c.Start()
}
