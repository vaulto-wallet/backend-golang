package handlers

import (
	"encoding/json"
	"net/http"
)

type JwtToken struct {
	Token string `json:"token"`
}

type Error int

const (
	Success Error = iota
	BadRequest
	NotImplemented
	NoUser
	IncorrectPassword
	AlreadyRegistered
	TokenMissing
	TokenInvalid
	TokenMalformed
)

func (error_code Error) String() string {
	error_text := [...]string{
		"SUCCESS",
		"BAD_REQUEST",
		"NOT_IMPLEMENTED",
		"NO_USER",
		"INCORRECT_PASSWORD",
		"ALREADY_REGISTERED",
		"AUTH_TOKEN_MISSING",
		"AUTH_TOKEN_INVALID",
		"AUTH_TOKEN_MALFORMED",
	}

	return error_text[error_code]
}

func (error_code Error) Text() string {
	error_text := [...]string{
		"Success",
		"Bad request",
		"Not implemented",
		"No such user",
		"Incorrect password",
		"Already registered",
		"Authorization token is missing",
		"Authorization token is invalid",
		"Authorization token is malformed",
	}
	return error_text[error_code]
}

type appError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
	Text  string `json:"text"`
}

type appBoolResult struct {
	Result bool `json:"result"`
}

type appStringResult struct {
	Result string `json:"result"`
}

type appJsonResult struct {
	Result interface{} `json:"result"`
}

func ReturnError(w http.ResponseWriter, e Error) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appError{e.String(), int(e), e.Text()})
}

func ReturnErrorWithStatus(w http.ResponseWriter, e Error, s int) {
	w.WriteHeader(s)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appError{e.String(), int(e), e.Text()})
}

func ReturnErrorWithStatusString(w http.ResponseWriter, e Error, s int, t string) {
	w.WriteHeader(s)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appError{e.String(), int(e), t})
}

func ReturnResult(w http.ResponseWriter, r interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appJsonResult{r})
}
