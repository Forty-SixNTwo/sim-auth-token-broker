module github.com/Forty-SixNTwo/sim-auth-token-broker-broker

go 1.24.3

require (
	github.com/Forty-SixNTwo/sim-auth-token-broker/libs/config v0.0.0
	github.com/Forty-SixNTwo/sim-auth-token-broker/libs/graceful v0.0.0
	github.com/Forty-SixNTwo/sim-auth-token-broker/libs/jwt v0.0.0
	github.com/Forty-SixNTwo/sim-auth-token-broker/libs/logs v0.0.0
	github.com/Forty-SixNTwo/sim-auth-token-broker/libs/utilities v0.0.0
	github.com/go-jose/go-jose/v4 v4.1.0
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/samber/slog-http v1.7.0 // indirect
	github.com/sony/gobreaker v1.0.0
	go.opentelemetry.io/otel v1.29.0 // indirect
	go.opentelemetry.io/otel/trace v1.29.0 // indirect
	golang.org/x/time v0.11.0
)

require (
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/sync v0.14.0
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/Forty-SixNTwo/sim-auth-token-broker/libs/config => ../../libs/config
	github.com/Forty-SixNTwo/sim-auth-token-broker/libs/graceful => ../../libs/graceful
	github.com/Forty-SixNTwo/sim-auth-token-broker/libs/jwt => ../../libs/jwt
	github.com/Forty-SixNTwo/sim-auth-token-broker/libs/logs => ../../libs/logs
	github.com/Forty-SixNTwo/sim-auth-token-broker/libs/utilities => ../../libs/utilities
)
