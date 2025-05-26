package utils

import (
	"testing"

	"github.com/Forty-SixNTwo/sim-auth-token-broker/libs/config"
)

func TestIsValidE164(t *testing.T) {
	cases := []struct {
		input string
		valid bool
	}{
		{"+972541234567", true},
		{"972541234567", true},
		{"+123", false},
		{"", false},
		{"+1 (234) 567-8901", false},
	}

	for _, c := range cases {
		got := IsValidE164(c.input)
		if got != c.valid {
			t.Errorf("IsValidE164(%q) = %v, want %v", c.input, got, c.valid)
		}
	}
}

func TestMatchPrefix(t *testing.T) {
	prefixMap := map[string]config.Telco{
		"97205":  {BaseURL: "orange"},
		"972050": {BaseURL: "vodafone"},
		"4477":   {BaseURL: "vf-uk"},
	}

	cases := []struct {
		phone   string
		wantURL string
		wantErr bool
	}{
		{"+97205012345", "vodafone", false},
		{"+97205123456", "orange", false},
		{"4477123456", "vf-uk", false},
		{"12345", "", true},
	}

	for _, c := range cases {
		telco, err := MatchPrefix(c.phone, prefixMap)
		if (err != nil) != c.wantErr {
			t.Errorf("MatchPrefix(%q) error = %v, wantErr %v", c.phone, err, c.wantErr)
			continue
		}
		if err == nil && telco.BaseURL != c.wantURL {
			t.Errorf("MatchPrefix(%q).BaseURL = %q, want %q", c.phone, telco.BaseURL, c.wantURL)
		}
	}
}
