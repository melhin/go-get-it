package middlewares

import (
	"net/http"

	"go-get-it/app/auth"
	"go-get-it/app/controllers"
)

func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			controllers.RespondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		next(w, r)
	}
}
