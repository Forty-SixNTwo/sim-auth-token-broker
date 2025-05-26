package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/Forty-SixNTwo/sim-auth-token-broker/libs/config"
	"github.com/Forty-SixNTwo/sim-auth-token-broker/libs/utilities"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/sony/gobreaker"
	"golang.org/x/sync/singleflight"
	"golang.org/x/time/rate"
)

const (
	jwksTTL   = 10 * time.Minute
	rateLimit = 5
	rateBurst = 10
)

var (
	jwksMu    sync.RWMutex
	jwksCache = make(map[string]*cachedSet)
	jwksGroup singleflight.Group
)

type cachedSet struct {
	set       jose.JSONWebKeySet
	fetchedAt time.Time
}

type TelcoClient struct {
	BaseURL      string
	ClientID     string
	ClientSecret string
	HTTP         *http.Client
	limiter      *rate.Limiter
	breaker      *gobreaker.CircuitBreaker
}

func New(cfgTelco config.Telco) *TelcoClient {
	cbSettings := gobreaker.Settings{
		Name:        cfgTelco.BaseURL,
		MaxRequests: 1,
		Interval:    time.Minute,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},
	}
	return &TelcoClient{
		BaseURL:      cfgTelco.BaseURL,
		ClientID:     cfgTelco.ClientID,
		ClientSecret: cfgTelco.ClientSecret,
		HTTP: &http.Client{
			Timeout:   5 * time.Second,
			Transport: utilities.RequestIDTransport(nil),
		},
		limiter: rate.NewLimiter(rate.Limit(rateLimit), rateBurst),
		breaker: gobreaker.NewCircuitBreaker(cbSettings),
	}
}

func (t *TelcoClient) ExchangeCode(ctx context.Context, form url.Values) (string, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := t.limiter.Wait(ctxWithTimeout); err != nil {
		return "", fmt.Errorf("rate limit wait failed: %w", err)
	}

	res, err := t.breaker.Execute(func() (any, error) {
		req, err := http.NewRequestWithContext(ctxWithTimeout, "POST", t.BaseURL+"/token", bytes.NewBufferString(form.Encode()))
		if err != nil {
			return "", err
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.SetBasicAuth(t.ClientID, t.ClientSecret)

		resp, err := t.HTTP.Do(req)
		if err != nil {
			return "", err
		}

		defer func() {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return "", fmt.Errorf("telco error %d: %s", resp.StatusCode, body)
		}

		var out struct {
			AccessToken string `json:"access_token"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			return "", err
		}
		return out.AccessToken, nil
	})

	if err != nil {
		return "", err
	}

	token, ok := res.(string)
	if !ok {
		return "", fmt.Errorf("unexpected response from circuit breaker")
	}
	return token, nil
}

func (t *TelcoClient) GetJWKs(ctx context.Context, jwksURL string) (jose.JSONWebKeySet, error) {
	jwksMu.RLock()
	if cs, ok := jwksCache[jwksURL]; ok && time.Since(cs.fetchedAt) < jwksTTL {
		set := cs.set
		jwksMu.RUnlock()
		return set, nil
	}
	jwksMu.RUnlock()

	v, err, _ := jwksGroup.Do(jwksURL, func() (any, error) {
		set, err := t.FetchJWKs(ctx, jwksURL)
		if err != nil {
			return jose.JSONWebKeySet{}, err
		}
		t.UpdateCache(jwksURL, set)
		return set, nil
	})
	if err != nil {
		return jose.JSONWebKeySet{}, err
	}
	return v.(jose.JSONWebKeySet), nil
}

func (t *TelcoClient) FetchJWKs(ctx context.Context, jwksURL string) (jose.JSONWebKeySet, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := t.limiter.Wait(ctxWithTimeout); err != nil {
		return jose.JSONWebKeySet{}, fmt.Errorf("rate limit wait failed: %w", err)
	}

	res, err := t.breaker.Execute(func() (any, error) {
		req, err := http.NewRequestWithContext(ctxWithTimeout, "GET", jwksURL, nil)
		if err != nil {
			return jose.JSONWebKeySet{}, fmt.Errorf("create jwks request: %w", err)
		}
		req.Header.Set("Accept", "application/json")
		resp, err := t.HTTP.Do(req)
		if err != nil {
			return jose.JSONWebKeySet{}, fmt.Errorf("fetch jwks: %w", err)
		}

		defer func() {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}()

		if resp.StatusCode != http.StatusOK {
			return jose.JSONWebKeySet{}, fmt.Errorf("jwks fetch status: %d", resp.StatusCode)
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return jose.JSONWebKeySet{}, fmt.Errorf("read jwks: %w", err)
		}
		var set jose.JSONWebKeySet
		if err := json.Unmarshal(data, &set); err != nil {
			return jose.JSONWebKeySet{}, fmt.Errorf("parse jwks: %w", err)
		}
		return set, nil
	})

	if err != nil {
		return jose.JSONWebKeySet{}, err
	}
	set, ok := res.(jose.JSONWebKeySet)
	if !ok {
		return jose.JSONWebKeySet{}, fmt.Errorf("unexpected response type from circuit breaker")
	}
	return set, nil
}

func (t *TelcoClient) UpdateCache(jwksURL string, set jose.JSONWebKeySet) {
	jwksMu.Lock()
	jwksCache[jwksURL] = &cachedSet{set: set, fetchedAt: time.Now()}
	jwksMu.Unlock()
}
