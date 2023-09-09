package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vadimistar/hackathon1/models"
)

type getVerification interface {
	GetVerification(token string) (*models.Verification, error)
}

type getUserByID interface {
	GetUserByID(id string) (*models.User, error)
}

func Verify(getVerify getVerification, getUser getUserByID, saveUser saveUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer respondOk(w, r)

		token := chi.URLParam(r, "token")
		if token == "" {
			log.Println("token is empty")
		}

		verification, err := getVerify.GetVerification(token)
		if err != nil {
			respondError(w, r, ErrInvalidToken)
			return
		}

		user, err := getUser.GetUserByID(verification.UserID)
		if err != nil {
			log.Printf("get user by id: %s", err)
			respondError(w, r, ErrInternalServer)
			return
		}

		// user.EmailVerified = 1

		err = saveUser.SaveUser(user)
		if err != nil {
			log.Printf("save user: %s", err)
			respondError(w, r, ErrInternalServer)
			return
		}

		w.WriteHeader(http.StatusOK)
		respondOk(w, r)
	}
}
