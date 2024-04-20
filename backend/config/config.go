package config

import "os"

var (
	IS_PRODUCTION = os.Getenv("DEPLOYMENT")
)
