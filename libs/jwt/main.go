package jwt

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-jose/go-jose/v4"
	Jwt "github.com/go-jose/go-jose/v4/jwt"

	"github.com/Forty-SixNTwo/sim-auth-token-broker-broker/clients"
	"github.com/Forty-SixNTwo/sim-auth-token-broker/libs/utilities"
)

const (
	ALG           = "alg"
	HS256         = "HS256"
	JWT           = "jwt"
	SIG           = "sig"
	CLIENT_ID     = "client_id"
	CLIENT_SECRET = "client_secret"
	GRANT_TYPE    = "grant_type"
	CODE          = "code"
	REALM         = `Basic realm="telco"`
)

var (
	privKey *rsa.PrivateKey
	pubJWK  jose.JSONWebKey
	signer  jose.Signer
)

type Payload struct {
	Issuer    string
	Subject   string
	Audience  []string
	ExpiresAt time.Time
	Extra     map[string]any
}

func Init(keyID string, bits int) error {
	var err error
	privKey, err = rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return fmt.Errorf("generate RSA key: %w", err)
	}

	pubJWK = jose.JSONWebKey{
		Key:       privKey.Public(),
		KeyID:     keyID,
		Algorithm: string(jose.RS256),
		Use:       SIG,
	}

	signer, err = jose.NewSigner(
		jose.SigningKey{Algorithm: jose.RS256, Key: privKey},
		(&jose.SignerOptions{}).WithType(JWT),
	)
	if err != nil {
		return fmt.Errorf("create signer: %w", err)
	}
	return nil
}

func InitHS256(secret []byte) error {
	var err error
	signer, err = jose.NewSigner(
		jose.SigningKey{Algorithm: jose.HS256, Key: secret},
		(&jose.SignerOptions{}).
			WithType(JWT).
			WithHeader(ALG, HS256),
	)
	return err
}

func JWKsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utilities.WriteJSONError(w, "method not allowed", r.Method, http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	jwks := jose.JSONWebKeySet{Keys: []jose.JSONWebKey{pubJWK}}
	json.NewEncoder(w).Encode(jwks)
}

func JWTsHandler(expectedID, expectedSecret, issuer string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utilities.WriteJSONError(w, "method not allowed", r.Method, http.StatusMethodNotAllowed)
			return
		}

		if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/x-www-form-urlencoded") {
			utilities.WriteJSONError(w, "unsupported media type", ct, http.StatusUnsupportedMediaType)
			return
		}

		if err := r.ParseForm(); err != nil {
			utilities.WriteJSONError(w, "invalid form", err.Error(), http.StatusBadRequest)
			return
		}

		id, secret, ok := r.BasicAuth()
		if !ok {
			id = r.PostFormValue(CLIENT_ID)
			secret = r.PostFormValue(CLIENT_SECRET)
		}
		if id != expectedID || secret != expectedSecret {
			w.Header().Set("WWW-Authenticate", REALM)
			utilities.WriteJSONError(w, "unauthorized", "", http.StatusUnauthorized)
			return
		}

		grantType := r.PostFormValue(GRANT_TYPE)
		code := r.PostFormValue(CODE)
		if grantType == "" || code == "" {
			utilities.WriteJSONError(w, "grant_type and code required", fmt.Sprintf("grantType %s code %s", grantType, code), http.StatusBadRequest)
			return
		}

		token, err := Sign(issuer, code, []string{id}, time.Hour)
		if err != nil {
			log.Printf("error signing token: %v", err)
			utilities.WriteJSONError(w, "internal error", err.Error(), http.StatusInternalServerError)
			return
		}

		resp := map[string]any{
			"access_token": token,
			"token_type":   "Bearer",
			"expires_in":   3600,
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(resp)
	}
}

func Sign(issuer, subject string, audience []string, expiresIn time.Duration) (string, error) {
	now := time.Now()
	claims := Jwt.Claims{
		Issuer:   issuer,
		Subject:  subject,
		Audience: Jwt.Audience(audience),
		IssuedAt: Jwt.NewNumericDate(now),
		Expiry:   Jwt.NewNumericDate(now.Add(expiresIn)),
	}
	return Jwt.Signed(signer).Claims(claims).Serialize()
}

func Validate(ctx context.Context, tokenStr string, tc *clients.TelcoClient, jwksURL string) (*Payload, error) {
	set, err := tc.GetJWKs(ctx, jwksURL)
	if err != nil {
		return nil, err
	}

	parsed, err := Jwt.ParseSigned(tokenStr, []jose.SignatureAlgorithm{jose.RS256})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	var claims Jwt.Claims
	for _, key := range set.Keys {
		if err := parsed.Claims(key.Key, &claims); err == nil {
			return &Payload{
				Issuer:    claims.Issuer,
				Subject:   claims.Subject,
				Audience:  claims.Audience,
				ExpiresAt: claims.Expiry.Time(),
			}, nil
		}
	}

	fresh, err := tc.FetchJWKs(ctx, jwksURL)
	if err != nil {
		return nil, fmt.Errorf("refresh jwks on kid miss: %w", err)
	}

	for _, key := range fresh.Keys {
		if err := parsed.Claims(key.Key, &claims); err == nil {
			return &Payload{
				Issuer:    claims.Issuer,
				Subject:   claims.Subject,
				Audience:  claims.Audience,
				ExpiresAt: claims.Expiry.Time(),
			}, nil
		}
	}

	return nil, fmt.Errorf("invalid token signature after refresh")
}

func Mint(p Payload) (string, error) {
	now := time.Now()
	std := Jwt.Claims{
		Issuer:   p.Issuer,
		Subject:  p.Subject,
		Audience: p.Audience,
		Expiry:   Jwt.NewNumericDate(p.ExpiresAt),
		IssuedAt: Jwt.NewNumericDate(now),
	}

	raw := struct {
		Jwt.Claims
		Extra map[string]any `json:"extra,omitempty"`
	}{
		Claims: std,
		Extra:  p.Extra,
	}
	tok, err := Jwt.Signed(signer).Claims(raw).Serialize()
	if err != nil {
		return "", fmt.Errorf("mint token: %w", err)
	}
	return tok, nil
}
