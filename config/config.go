package config

import (
	"os"
)

var ApiKey = getEnv("API_KEY", "your-secure-api-key")

const BrokerAddress = "csqnjbdjp6ucv9qq9fi0.any.us-east-1.mpx.prd.cloud.redpanda.com:9092"

const RateLimitRequestsPerSecond = 5
const RateLimitBurst = 10

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
