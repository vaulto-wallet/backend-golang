package middlewares

import (
	h "../handlers"
	m "../models"
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strings"
)

func AuthMiddlewareGenerator(db *gorm.DB) (mw func(http.Handler) http.Handler) {
	mw = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.RequestURI)

			notAuth := []string{"/api/users/login", "/api/users/register", "/api/clear", "/api/start"}
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
				h.ReturnErrorWithStatus(w, h.TokenMissing, http.StatusForbidden)
				return
			}

			splitted := strings.Split(tokenHeader, " ")
			if len(splitted) != 2 {
				h.ReturnErrorWithStatus(w, h.TokenInvalid, http.StatusForbidden)
				return
			}

			// obtain JWT token
			tokenPart := splitted[1]
			tk := h.AuthToken{}

			// parse JWT token
			token, err := jwt.ParseWithClaims(tokenPart, &tk, func(token *jwt.Token) (interface{}, error) {
				return []byte("cryptosecret"), nil
			})

			// cannot parse JWT token
			if err != nil {
				h.ReturnErrorWithStatus(w, h.TokenMalformed, http.StatusForbidden)
				return
			}

			// token is not valid
			if !token.Valid {
				h.ReturnErrorWithStatus(w, h.TokenInvalid, http.StatusForbidden)
				return
			}

			dbUser := new(m.User)
			db.Preload("Account").First(&dbUser, "username = ?", tk.Username)

			ctx := context.WithValue(r.Context(), "user", dbUser)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
	return
}
