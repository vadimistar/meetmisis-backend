package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
)

func userIDFromCookie(w http.ResponseWriter, r *http.Request, jwtKey []byte) string {
	tokenCookie, err := r.Cookie("meetmisis-token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			respondError(w, r, ErrNoAuth)
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

	if userID, ok := claims["userID"]; ok {
		if userID != "" {
			return claims["userID"].(string)
		}
	}

	respondError(w, r, ErrNoAuth)
	return ""
}
