module github.com/Forty-SixNTwo/sim-auth-token-broker-partner

go 1.24.3

require (
	github.com/Forty-SixNTwo/sim-auth-token-broker/libs/config v0.0.0
	github.com/Forty-SixNTwo/sim-auth-token-broker/libs/graceful v0.0.0
	github.com/Forty-SixNTwo/sim-auth-token-broker/libs/jwt v0.0.0
	github.com/Forty-SixNTwo/sim-auth-token-broker/libs/logs v0.0.0
)

require (
	github.com/go-jose/go-jose/v4 v4.1.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/Forty-SixNTwo/sim-auth-token-broker/libs/config => ../../libs/config
	github.com/Forty-SixNTwo/sim-auth-token-broker/libs/graceful => ../../libs/graceful
	github.com/Forty-SixNTwo/sim-auth-token-broker/libs/jwt => ../../libs/jwt
	github.com/Forty-SixNTwo/sim-auth-token-broker/libs/logs => ../../libs/logs
)
