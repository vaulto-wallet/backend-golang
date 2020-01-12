package middlewares

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"os"
	"strings"
	h "../handlers"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)

		notAuth := []string{"/api/users/login", "/api/users/register", "/api/clear"}
		requestPath := r.URL.Path

		// check if path doesn't require authorization
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("Authorization") // get token from HTTP header

		if tokenHeader == "" {
			h.ThrowErrorWithStatus(w, h.TokenMissing, http.StatusForbidden)
			return
		}


		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			h.ThrowErrorWithStatus(w, h.TokenInvalid, http.StatusForbidden)
			return
		}

		// obtain JWT token
		tokenPart := splitted[1]
		tk := h.AuthToken{}

		// parse JWT token
		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		// cannot parse JWT token
		if err != nil {
			h.ThrowErrorWithStatus(w, h.TokenMalformed, http.StatusForbidden)
			return
		}

		// token is not valid
		if !token.Valid {
			h.ThrowErrorWithStatus(w, h.TokenInvalid, http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "user", tk.UserName)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
