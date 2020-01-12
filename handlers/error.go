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
	NoUser
	IncorrectPassword
	AlreadyRegistered
	TokenMissing
	TokenInvalid
	TokenMalformed
)

func (error_code Error) String() string {
	error_text  := [...]string{
		"SUCCESS",
		"BAD_REQUEST",
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
	error_text  := [...]string{
		"Success",
		"Bad request",
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
	Code int `json:"code"`
	Text string `json:"text"`
}

type appBoolResult struct {
	Result bool `json:"result"`
}

type appStringResult struct {
	Result string `json:"result"`
}


func ThrowError(w http.ResponseWriter, e Error)  {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appError{e.String() , int(e), e.Text()  })
}

func ThrowErrorWithStatus(w http.ResponseWriter, e Error, s int)  {
	w.WriteHeader(s)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appError{e.String() , int(e), e.Text()  })
}

func ReturnBool(w http.ResponseWriter, r bool)  {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appBoolResult{ r })
}

func ReturnString(w http.ResponseWriter, r string)  {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appStringResult{ r })
}
