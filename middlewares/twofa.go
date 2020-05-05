package middlewares

import (
	h "../handlers"
	m "../models"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

func TwoFAMiddlewareGenerator(db *gorm.DB) (mw func(http.Handler) http.Handler) {
	mw = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			required := []struct {
				Path   string
				Method string
			}{{"/api/orders/confirm", "POST"}}

			requestPath := r.URL.Path
			requestMethod := r.Method

			masterPassword := r.Context().Value("OTP")
			// check if path doesn't require engine was started
			for _, value := range required {
				if value.Path == requestPath && value.Method == requestMethod && masterPassword == nil {
					user := r.Context().Value("user").(m.User)
					log.Println("Two FA Middleware", r.RequestURI, user.Username)

					if user.Account.OTPStatus != m.OTPStatusActive {
						h.ReturnErrorWithStatusString(w, h.TokenInvalid, http.StatusForbidden, "OTP Required")
						return
					}

					OTPHeader := r.Header.Get("OTPAuthorization") // get token from HTTP header
					if !user.Account.VerifyTOTP(OTPHeader) {
						h.ReturnErrorWithStatusString(w, h.TokenInvalid, http.StatusForbidden, "Invalid OTP")
						return
					}
				}
			}

			next.ServeHTTP(w, r)
			return

		})
	}
	return
}
