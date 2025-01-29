package config

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

const OIDC_PROVIDER_INITILIZATION_MAX_RETRIES = 5
const JWTSourceTimeout = 15 * time.Second

type spiffeIDs struct {
	Tokenetes spiffeid.ID
	Order     spiffeid.ID
	Stocks    spiffeid.ID
}

type GatewayConfig struct {
	TokenetesURL     *url.URL
	StocksServiceURL *url.URL
	OrderServiceURL  *url.URL
	SpiffeIDs        *spiffeIDs
	TraTToggle       bool
}

func GetAppConfig() *GatewayConfig {
	return &GatewayConfig{
		TokenetesURL:     parseURL(getEnv("TOKENETES_URL")),
		StocksServiceURL: parseURL(getEnv("STOCKS_SERVICE_URL")),
		OrderServiceURL:  parseURL(getEnv("ORDER_SERVICE_URL")),
		SpiffeIDs: &spiffeIDs{
			Tokenetes: spiffeid.RequireFromString(getEnv("TOKENETES_SPIFFE_ID")),
			Order:     spiffeid.RequireFromString(getEnv("ORDER_SERVICE_SPIFFE_ID")),
			Stocks:    spiffeid.RequireFromString(getEnv("STOCKS_SERVICE_SPIFFE_ID")),
		},
		TraTToggle: getBoolEnv("ENABLE_TRATS"),
	}
}

func GetOauth2Config() oauth2.Config {
	return oauth2.Config{
		ClientID:     getEnv("OAUTH2_CLIENT_ID"),
		ClientSecret: getEnv("OAUTH2_CLIENT_SECRET"),
		RedirectURL:  getEnv("OAUTH2_REDIRECT_URL"),
		Endpoint: oauth2.Endpoint{
			TokenURL: getEnv("OAUTH2_TOKEN_URL"),
		},
		Scopes: []string{"openid", "profile", "email"},
	}
}

func GetOIDCProvider(logger *zap.Logger) *oidc.Provider {
	delay := time.Second

	for i := 0; i < OIDC_PROVIDER_INITILIZATION_MAX_RETRIES; i++ {
		ctx := context.Background()
		oidcIssuer := getEnv("OIDC_ISSUER_URL")

		provider, err := oidc.NewProvider(ctx, oidcIssuer)
		if err == nil {
			logger.Info("Successfully connected to the OIDC provider.")

			return provider
		}

		logger.Error("Failed to connect to the OIDC provider.",
			zap.Int("attempt", i+1),
			zap.String("retrying_in", delay.String()),
			zap.Error(err))
		time.Sleep(delay)

		delay *= 2
	}

	logger.Error(fmt.Sprintf("Failed to connect to the OIDC provider after %d attempts", OIDC_PROVIDER_INITILIZATION_MAX_RETRIES))

	panic(fmt.Sprintf("failed to connect to the OIDC provider after %d attempts", OIDC_PROVIDER_INITILIZATION_MAX_RETRIES))
}

func getBoolEnv(key string) bool {
	val, err := strconv.ParseBool(getEnv(key))
	if err != nil {
		panic("Error parsing boolean environment variable " + key + ": " + err.Error())
	}

	return val
}

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		panic(fmt.Sprintf("%s environment variable not set", key))
	}

	return value
}

func parseURL(rawurl string) *url.URL {
	parsedURL, err := url.Parse(rawurl)
	if err != nil {
		panic(fmt.Sprintf("Error parsing URL %s: %v", rawurl, err))
	}

	return parsedURL
}
