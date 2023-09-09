package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
)

func userIDFromCookie(w http.ResponseWriter, r *http.Request, noCookieRedirectURL string, jwtKey []byte) string {
	tokenCookie, err := r.Cookie("meetmisis-token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			http.Redirect(w, r, noCookieRedirectURL, http.StatusMovedPermanently)
			return ""
		}

		log.Printf("get cookie: %s", err)
		respondError(w, r, ErrInternalServer)
		return ""
	}

	claims := make(jwt.MapClaims)
	_, err = jwt.ParseWithClaims(tokenCookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		log.Printf("cannot parse jwt: %s", err)
		respondError(w, r, ErrInternalServer)
		return ""
	}

	return claims["userID"].(string)
}
