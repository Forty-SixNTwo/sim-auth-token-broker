package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Forty-SixNTwo/sim-auth-token-broker-broker/service"
	"github.com/Forty-SixNTwo/sim-auth-token-broker/libs/config"
	"github.com/Forty-SixNTwo/sim-auth-token-broker/libs/graceful"
	"github.com/Forty-SixNTwo/sim-auth-token-broker/libs/jwt"
	"github.com/Forty-SixNTwo/sim-auth-token-broker/libs/logs"
)

func main() {
	logger := logs.Init("[Broker]")

	cfg, err := config.LoadBrokerConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	if err := jwt.InitHS256([]byte(cfg.SigningKey)); err != nil {
		logs.Fatal(logger, "jwt init failed", "error", err)
	}

	mux := http.NewServeMux()
	handler := service.NewTokenHandler(cfg, logger)
	mux.Handle("/token",
		logs.LoggingMiddleware(logger)(
			http.HandlerFunc(handler.Handle),
		),
	)

	srv := &http.Server{
		Addr:    cfg.ListenAddr,
		Handler: mux,
	}

	if err := graceful.StartServer(srv, 5*time.Second, logger); err != nil {
		logs.Fatal(logger, "server failure", "error", err)
	}
}
