package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const formMaxMemory = 4 * 1024 * 1024

func Login(getUser getUser, jwtKey []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(formMaxMemory)
		if err != nil {
			log.Printf("cannot parse form: %s", err)
			respondError(w, r, ErrInternalServer)
			return
		}

		// maybe validation

		user, err := getUser.GetUser(r.FormValue("username"))
		if err != nil {
			log.Printf("cannot get user: %s", err)
			respondError(w, r, ErrInternalServer)
			return
		}

		if user == nil || user.ID == "" {
			respondError(w, r, ErrNoUserWithEmail)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.FormValue("password")))
		if err != nil {
			respondError(w, r, ErrInvalidCredentials)
			return
		}

		tokenString, err := authToken(user.ID, jwtKey)
		if err != nil {
			respondError(w, r, ErrInternalServer)
			log.Printf("error while generate jwt token: %s", err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "meetmisis-token",
			Value:    tokenString,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, successResponse{
			Response: Response{
				Status: statusSuccess,
			},
			Data: registerLoginResponse{
				Token: tokenString,
			},
		})
	}
}

type registerLoginResponse struct {
	Token string `json:"token"`
}

func authToken(userID string, jwtKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    userID,
		"createdAt": time.Now(),
	})
	return token.SignedString(jwtKey)
}
