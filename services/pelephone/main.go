package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/Forty-SixNTwo/sim-auth-token-broker/libs/config"
	"github.com/Forty-SixNTwo/sim-auth-token-broker/libs/graceful"
	"github.com/Forty-SixNTwo/sim-auth-token-broker/libs/jwt"
	"github.com/Forty-SixNTwo/sim-auth-token-broker/libs/logs"
)

const (
	TELCO_KEY_ID        = "PELEPHONE_KEY_ID"
	TELCO_CLIENT_ID     = "PELEPHONE_CLIENT_ID"
	TELCO_CLIENT_SECRET = "PELEPHONE_CLIENT_SECRET"
	TELCO_ISSUER_URL    = "PELEPHONE_ISSUER_URL"
)

func main() {
	logger := logs.Init("[Pelephone-Telco]")
	cfg, err := config.LoadTelcoConfig(TELCO_KEY_ID, TELCO_CLIENT_ID, TELCO_CLIENT_SECRET, TELCO_ISSUER_URL)
	if err != nil {
		logs.Fatal(logger, "config init failed", "error", err)
	}

	port := strings.TrimPrefix(cfg.TelcoIssuerURL, "http://localhost")

	if err := jwt.Init(cfg.TelcoKeyID, 2048); err != nil {
		logs.Fatal(logger, "jwt init failed", "error", err)
	}

	mux := http.NewServeMux()
	mux.Handle(
		"/.well-known/jwks.json",
		logs.LoggingMiddleware(logger)(
			http.HandlerFunc(jwt.JWKsHandler),
		),
	)
	mux.Handle(
		"/token",
		logs.LoggingMiddleware(logger)(
			http.HandlerFunc(jwt.JWTsHandler(cfg.TelcoClientID, cfg.TelcoClientSecret, cfg.TelcoIssuerURL)),
		),
	)

	srv := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	if err := graceful.StartServer(srv, 5*time.Second, logger); err != nil {
		logs.Fatal(logger, "server failure", "error", err)
	}
}
