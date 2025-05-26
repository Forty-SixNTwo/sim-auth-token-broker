package model

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestParse_Success(t *testing.T) {
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {"abc123"},
		"phone":         {"+972541234567"},
		"redirect_uri":  {"https://example.com/callback"},
		"code_verifier": {"verifier"},
	}

	req, err := http.NewRequest(http.MethodPost, "/token", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	tr, err := Parse(req)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if tr.GrantType != "authorization_code" ||
		tr.Code != "abc123" ||
		tr.Phone != "+972541234567" ||
		tr.RedirectURI != "https://example.com/callback" ||
		tr.CodeVerifier != "verifier" {
		t.Errorf("Parse returned %+v; want values from form", tr)
	}
}

func TestParse_BadForm(t *testing.T) {
	req := &http.Request{Body: http.NoBody}
	_, err := Parse(req)
	if err == nil {
		t.Error("Parse expected error for invalid form")
	}
}
