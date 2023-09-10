package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"github.com/vadimistar/hackathon1/models"
)

type getPartner interface {
	GetPartner(userId string) (*models.Partner, error)
}

type savePartner interface {
	SavePartner(p *models.Partner) error
}

type getUsersWithTags interface {
	GetUsersWithTags(tags []string) ([]string, error)
}

func FindMatch(p getPartner, sp savePartner, gu getUsersWithTags, gt getTaggedUser, jwtKey []byte) http.HandlerFunc {
	var h http.HandlerFunc

	h = func(w http.ResponseWriter, r *http.Request) {
		userID := userIDFromCookie(w, r, jwtKey)
		if userID == "" {
			return
		}

		defer r.Body.Close()

		var req findMatchRequest

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			respondError(w, r, ErrInvalidInputData)
			return
		}

		partner, err := p.GetPartner(userID)
		log.Printf("%+v", partner)
		if err != nil {
			log.Printf("cannot get user by id (id = %s): err = %s", userID, err)
			respondError(w, r, ErrInternalServer)
			return
		}
		if (partner == nil || partner.Partners == nil) && err == nil {
			err := sp.SavePartner(&models.Partner{
				UserID:   userID,
				Partners: []string{},
			})
			if err != nil {
				log.Printf("cannot create partner for userID = %s: %s", userID, err)
				respondError(w, r, ErrInternalServer)
				return
			}
		}

		partners := make(map[string]struct{})
		for _, p := range partner.Partners {
			partners[p] = struct{}{}
		}

		matchedUsers, err := matchUsersByTags(gu, req.Tags)
		if err != nil {
			log.Printf("match users by tags: %s", err.Error())
			respondError(w, r, ErrInternalServer)
			return
		}

		for _, matchedUser := range matchedUsers {
			if _, ok := partners[matchedUser]; ok {
				continue
			}

			err := sp.SavePartner(&models.Partner{
				UserID:   userID,
				Partners: append(partner.Partners, matchedUser),
			})
			if err != nil {
				log.Printf("add partner %s: %s", userID, err)
				respondError(w, r, ErrInternalServer)
				return
			}

			render.JSON(w, r, successResponse{
				Response: Response{
					Status: statusSuccess,
				},
				Data: findMatchResponse{
					Id: matchedUser,
				},
			})
			return
		}

		err = sp.SavePartner(&models.Partner{
			UserID:   userID,
			Partners: []string{},
		})
		if err != nil {
			log.Printf("remove all partners %s: %s", userID, err)
			respondError(w, r, ErrInternalServer)
			return
		}

		h(w, r)
	}

	return h
}

func matchUsersByTags(gu getUsersWithTags, tags []string) ([]string, error) {
	matchedUsers, err := gu.GetUsersWithTags(tags)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get users with tags %+v", tags)
	}

	// todo: maybe add cache here to not access db every single time again

	if len(matchedUsers) == 0 {
		if len(tags) > 0 {
			return matchUsersByTags(gu, tags[1:])
		}

		// has to unreachable, because if len(tags) == 0, then all users do match,
		// there are no constrains
		return nil, errors.New("unreachable: all tags are matched")
	}

	return matchedUsers, nil
}

type findMatchRequest struct {
	Tags []string `json:"tags"`
}

type findMatchResponse struct {
	Id string `json:"id"`
}
