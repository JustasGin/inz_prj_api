package main

import (
	"errors"
	"github.com/pascaldekloe/jwt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		next.ServeHTTP(w, r)
	})
}

func (app *application) checkToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authHeader := r.Header.Get("Authorization")

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			app.errorJSON(w, errors.New("Invalid auth header"))
			return
		}
		if headerParts[0] != "Bearer" {
			app.errorJSON(w, errors.New("No bearer"))
			return
		}

		token := headerParts[1]
		claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.jwt.secret))
		if err != nil {
			app.errorJSON(w, errors.New("Unauthorized - failed HMAC check"), http.StatusForbidden)
			return
		}
		if !claims.Valid(time.Now()) {
			app.errorJSON(w, errors.New("Unauthorized - token invalid"), http.StatusForbidden)
			return
		}
		if !claims.AcceptAudience(os.Getenv("JWT_DNS")) {
			app.errorJSON(w, errors.New("Unauthorized - invalid audience"), http.StatusForbidden)
			return
		}
		if claims.Issuer != os.Getenv("JWT_DNS") {
			app.errorJSON(w, errors.New("Unauthorized - invalid issuer"), http.StatusForbidden)
			return
		}

		_, err = strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			app.errorJSON(w, errors.New("Unauthorized"), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
