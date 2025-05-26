package model

import "net/http"

type TokenRequest struct {
	GrantType    string
	Code         string
	Phone        string
	RedirectURI  string
	CodeVerifier string
}

func Parse(r *http.Request) (TokenRequest, error) {
	if err := r.ParseForm(); err != nil {
		return TokenRequest{}, err
	}
	return TokenRequest{
		GrantType:    r.PostFormValue("grant_type"),
		Code:         r.PostFormValue("code"),
		Phone:        r.PostFormValue("phone"),
		RedirectURI:  r.PostFormValue("redirect_uri"),
		CodeVerifier: r.PostFormValue("code_verifier"),
	}, nil
}
