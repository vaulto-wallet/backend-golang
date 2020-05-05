package middlewares

import (
	h "../handlers"
	"log"
	"net/http"
)

func StartedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)

		required := []struct {
			Path   string
			Method string
		}{{"/api/seeds", "POST"}, {"/api/wallets", "POST"}, {"/api/address", "POST"}}

		requestPath := r.URL.Path
		requestMethod := r.Method
		masterPassword := r.Context().Value("masterPassword")
		// check if path doesn't require engine was started
		for _, value := range required {
			if value.Path == requestPath && value.Method == requestMethod && masterPassword == nil {
				h.ReturnErrorWithStatusString(w, h.BadRequest, http.StatusForbidden, "Encryption engine is not started")
				return
			}
		}

		next.ServeHTTP(w, r)
		return

	})
}
