package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/render"
	"github.com/vadimistar/hackathon1/models"
	"golang.org/x/crypto/bcrypt"
)

// type getUser interface {
// 	GetUser(username string, email string) (*models.User, error)
// }

type getUser interface {
	GetUser(username string) (*models.User, error)
}

type saveUser interface {
	SaveUser(user *models.User) error
}

type saveVerification interface {
	SaveVerification(v *models.Verification) error
}

// type EmailCredentials struct {
// 	Username string
// 	Password string
// }

func Register(getUser getUser, saveUser saveUser, saveVerify saveVerification, jwtKey []byte /*, endpoint string */ /* emailCredentials EmailCredentials */) (http.HandlerFunc, error) {
	// _, err := url.Parse(endpoint)
	// if err != nil {
	// 	return nil, fmt.Errorf("cannot parse endpoint, that is provided: %s", endpoint)
	// }

	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(formMaxMemory)
		if err != nil {
			log.Printf("cannot parse form: %s", err)
			respondError(w, r, ErrInternalServer)
			return
		}

		// TODO ADD VALIDATION
		username := r.FormValue("username")

		user, err := getUser.GetUser(username)
		if err != nil {
			log.Printf("cannot get user: %s", err)
			respondError(w, r, ErrInternalServer)
			return
		}

		if user != nil {
			if user.Username == username {
				respondError(w, r, ErrSameNicknameExists)
				return
			}
		}

		password := r.FormValue("password")
		// validate password!

		hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("bcrypt generate from password: %s", err)
			respondError(w, r, ErrInternalServer)
			return
		}

		newUser := &models.User{
			Username: username,
			Password: string(hashedPass),
		}

		err = saveUser.SaveUser(newUser)
		if err != nil {
			log.Printf("cannot save user: %s", err)
			respondError(w, r, ErrInternalServer)
			return
		}

		// verificationToken, err := generateToken(32)
		// if err != nil {
		// 	log.Printf("cannot generate verification token: %s", err)
		// 	respondError(w, r, ErrInternalServer)
		// 	return
		// }

		// verification := &models.Verification{
		// 	UserID: newUser.ID,
		// 	Token:  verificationToken,
		// }

		// if err := saveVerify.SaveVerification(verification); err != nil {
		// 	log.Printf("cannot save verification: %s", err)
		// 	respondError(w, r, ErrInternalServer)
		// 	return
		// }

		// todo: add email verification

		// 		err = mailru.Send(emailCredentials.Username,
		// 			emailCredentials.Password,
		// 			[]string{newUser.Email},
		// 			[]byte(fmt.Sprintf("Subject: MeetMisis: Регистрация\r\n\r\n"+`Вы зарегистрировались на платформе MeetMisis!

		// Подтвердите свою электронную почту:
		// %s/verify?token=%s
		// `, endpoint, verificationToken)))
		// 		if err != nil {
		// 			log.Printf("cannot send verification email to address: %s (%s)", newUser.Email, err)
		// 			respondError(w, r, ErrInternalServer)
		// 			return
		// 		}

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

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, registerLoginResponse{
			Response: Response{
				Status: statusSuccess,
			},
			Token: tokenString,
		})
	}, nil
}

// func generateToken(length int) (string, error) {
// 	token := make([]byte, length)
// 	_, err := rand.Read(token)
// 	if err != nil {
// 		return "", err
// 	}
// 	return base64.URLEncoding.EncodeToString(token), nil
// }
