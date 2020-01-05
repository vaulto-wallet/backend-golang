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
		// Do stuff here
		log.Println(r.RequestURI)

		notAuth := []string{"/api/user/login"} //Список эндпоинтов, для которых не требуется авторизация
		requestPath := r.URL.Path //текущий путь запроса

		//проверяем, не требует ли запрос аутентификации, обслуживаем запрос, если он не нужен
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("Authorization") //Получение токена

		if tokenHeader == "" { //Токен отсутствует, возвращаем  403 http-код Unauthorized
			h.ThrowErrorWithStatus(w, h.TokenMissing, http.StatusForbidden)
			return
		}

		splitted := strings.Split(tokenHeader, " ") //Токен обычно поставляется в формате `Bearer {token-body}`, мы проверяем, соответствует ли полученный токен этому требованию
		if len(splitted) != 2 {
			h.ThrowErrorWithStatus(w, h.TokenInvalid, http.StatusForbidden)
			return
		}

		tokenPart := splitted[1] //Получаем вторую часть токена
		tk := h.AuthToken{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil { //Неправильный токен, как правило, возвращает 403 http-код
			h.ThrowErrorWithStatus(w, h.TokenMalformed, http.StatusForbidden)
			return
		}

		if !token.Valid { //токен недействителен, возможно, не подписан на этом сервере
			h.ThrowErrorWithStatus(w, h.TokenInvalid, http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "user", tk.UserName)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
