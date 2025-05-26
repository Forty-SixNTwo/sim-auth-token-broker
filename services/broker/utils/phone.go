package utils

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/Forty-SixNTwo/sim-auth-token-broker/libs/config"
)

var e164Regex = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)

func IsValidE164(phone string) bool {
	return e164Regex.MatchString(phone)
}

func MatchPrefix(phone string, prefixMap map[string]config.Telco) (config.Telco, error) {
	pn := strings.TrimPrefix(strings.TrimSpace(phone), "+")
	keys := make([]string, 0, len(prefixMap))
	for p := range prefixMap {
		keys = append(keys, p)
	}
	sort.Slice(keys, func(i, j int) bool { return len(keys[i]) > len(keys[j]) })
	for _, p := range keys {
		if strings.HasPrefix(pn, p) {
			return prefixMap[p], nil
		}
	}
	return config.Telco{}, fmt.Errorf("no prefix match for %q", pn)
}
