package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Credentials struct {
	Username string `json:"email" bson:"email"`
	Password string `json:"password" json:"password"`
}

func (app *application) signIn(w http.ResponseWriter, r *http.Request) {
	var cred Credentials

	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		app.errorJSON(w, errors.New("Unauthorized access"))
		return
	}

	user, err := app.models.DB.GetUser(cred.Username)
	if err != nil {
		app.errorJSON(w, errors.New("Unauthorized access"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cred.Password))
	if err != nil {
		app.errorJSON(w, errors.New("Unauthorized access"))
		return
	}

	var claims jwt.Claims
	claims.Subject = fmt.Sprint(user.ID)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(2 * time.Hour))
	claims.Issuer = os.Getenv("JWT_DNS")
	claims.Audiences = []string{os.Getenv("JWT_DNS")}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secret))
	if err != nil {
		app.errorJSON(w, errors.New("error signing"))
		return
	}

	app.writeJSON(w, http.StatusOK, string(jwtBytes), "response")
}

func (app *application) checkUp(w http.ResponseWriter, r *http.Request) {
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

	app.writeJSON(w, http.StatusOK, "", "response")
}
