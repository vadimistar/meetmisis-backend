package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/render"
	"github.com/vadimistar/hackathon1/models"
)

type getTaggedUser interface {
	GetTaggedUser(userID string) (*models.TaggedUser, error)
}

type getTag interface {
	GetTag(tagID string) (*models.Tag, error)
}

func GetTags(gu getTaggedUser, gt getTag, jwtKey []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := userIDFromCookie(w, r, jwtKey)
		if userID == "" {
			return
		}

		taggedUser, err := gu.GetTaggedUser(userID)
		if err != nil {
			log.Printf("get tagged user with id=%s (%s)", userID, err)
			respondError(w, r, ErrInternalServer)
			return
		}

		if len(taggedUser.Tags) == 0 {
			render.JSON(w, r, successResponse{
				Response: Response{
					Status: statusSuccess,
				},
				Data: []models.Tag{},
			})
			return
		}

		tags := []*models.Tag{}
		for _, tagId := range taggedUser.Tags {
			tag, err := gt.GetTag(tagId)
			if err != nil {
				log.Printf("cannot get tag for user=%s id=%s: %s", taggedUser.UserID, tagId, err)
				respondError(w, r, ErrInvalidInputData)
				return
			}
			tags = append(tags, tag)
		}

		render.JSON(w, r, successResponse{
			Response: Response{
				Status: statusSuccess,
			},
			Data: tags,
		})
	}
}
