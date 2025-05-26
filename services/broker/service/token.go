package service

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Forty-SixNTwo/sim-auth-token-broker-broker/clients"
	"github.com/Forty-SixNTwo/sim-auth-token-broker-broker/model"
	"github.com/Forty-SixNTwo/sim-auth-token-broker-broker/utils"
	"github.com/Forty-SixNTwo/sim-auth-token-broker/libs/config"
	"github.com/Forty-SixNTwo/sim-auth-token-broker/libs/jwt"
	"github.com/Forty-SixNTwo/sim-auth-token-broker/libs/utilities"
)

type TokenHandler struct {
	cfg    *config.BrokerConfig
	logger *slog.Logger
}

func NewTokenHandler(cfg *config.BrokerConfig, logger *slog.Logger) *TokenHandler {
	return &TokenHandler{cfg: cfg, logger: logger}
}

func (h *TokenHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utilities.WriteJSONError(w, "method not allowed", r.Method, http.StatusMethodNotAllowed)
		return
	}
	if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/x-www-form-urlencoded") {
		utilities.WriteJSONError(w, "unsupported media type", ct, http.StatusUnsupportedMediaType)
		return
	}

	req, err := model.Parse(r)
	if err != nil {
		utilities.WriteJSONError(w, "invalid_form", err.Error(), http.StatusBadRequest)
		return
	}
	if req.GrantType != "authorization_code" {
		utilities.WriteJSONError(w, "unsupported_grant_type", "only authorization_code is supported", http.StatusBadRequest)
		return
	}

	if !utils.IsValidE164(req.Phone) {
		utilities.WriteJSONError(w, "invalid phone number format", "only E.164 phone numbers are supported", http.StatusBadRequest)
		return
	}

	telcoCfg, err := utils.MatchPrefix(req.Phone, h.cfg.PrefixMap)
	if err != nil {
		utilities.WriteJSONError(w, "invalid phone number", err.Error(), http.StatusBadRequest)
		return
	}

	tel := clients.New(telcoCfg)
	form := url.Values{
		"grant_type":    {req.GrantType},
		"code":          {req.Code},
		"redirect_uri":  {req.RedirectURI},
		"code_verifier": {req.CodeVerifier},
	}
	access, err := tel.ExchangeCode(r.Context(), form)
	if err != nil {
		utilities.WriteJSONError(w, "unable to retrive", err.Error(), http.StatusBadGateway)
		return
	}

	claims, err := jwt.Validate(r.Context(), access, tel, telcoCfg.BaseURL+"/.well-known/jwks.json")
	if err != nil {
		utilities.WriteJSONError(w, "invalid token from telco", err.Error(), http.StatusBadGateway)
		return
	}

	outToken, err := jwt.Mint(jwt.Payload{
		Issuer:    "sim-broker",
		Subject:   claims.Subject,
		Audience:  claims.Audience,
		ExpiresAt: time.Now().Add(15 * time.Minute),
		Extra:     map[string]any{"auth_method": "sim"},
	})
	if err != nil {
		utilities.WriteJSONError(w, "cannot mint token", err.Error(), http.StatusInternalServerError)
		return
	}

	resp := model.TokenResponse{
		AccessToken: outToken,
		TokenType:   "bearer",
		ExpiresIn:   900,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
