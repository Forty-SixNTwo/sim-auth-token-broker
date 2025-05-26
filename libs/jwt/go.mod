module github.com/Forty-SixNTwo/sim-auth-token-broker/libs/jwt

go 1.24.3

require github.com/go-jose/go-jose/v4 v4.1.0

require (
	github.com/Forty-SixNTwo/sim-auth-token-broker/libs/utilities v0.0.0
)

replace github.com/Forty-SixNTwo/sim-auth-token-broker/libs/utilities => ../../libs/utilities
