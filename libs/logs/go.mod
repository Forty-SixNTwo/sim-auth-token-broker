module github.com/Forty-SixNTwo/sim-auth-token-broker/libs/logs

go 1.24.3

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/samber/slog-http v1.7.0
	go.opentelemetry.io/otel v1.29.0 // indirect
	go.opentelemetry.io/otel/trace v1.29.0 // indirect
)

require github.com/Forty-SixNTwo/sim-auth-token-broker/libs/utilities v0.0.0

replace github.com/Forty-SixNTwo/sim-auth-token-broker/libs/utilities => ../utilities
