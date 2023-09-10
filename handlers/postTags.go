package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/render"
	"github.com/vadimistar/hackathon1/models"
)

type saveTagsForUser interface {
	SaveTagsForUser(userID string, tags []string) error
}

type saveTag interface {
	SaveTag(tag *models.Tag) error
}

func PostTags(st saveTag, stu saveTagsForUser, jwtKey []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := userIDFromCookie(w, r, jwtKey)
		if userID == "" {
			return
		}

		defer r.Body.Close()

		var request postTagsRequest
		err := render.DecodeJSON(r.Body, &request)
		if err != nil {
			log.Printf("decode json for user=%s: %s", userID, err)
			respondError(w, r, ErrInvalidInputData)
			return
		}

		for _, tag := range request.Tags {
			if tag.Id == "" {
				// TODO: Refactor saveTag to save multiple tags at once
				err := st.SaveTag(tag)
				if err != nil {
					log.Printf("cannot save tag: %+v (%s)", tag, err)
					respondError(w, r, ErrInternalServer)
					return
				}
			}
		}

		var tagIds []string
		for _, tag := range request.Tags {
			tagIds = append(tagIds, tag.Id)
		}

		err = stu.SaveTagsForUser(userID, tagIds)
		if err != nil {
			log.Printf("save tags: %s", err)
			respondError(w, r, ErrInvalidInputData)
			return
		}

		respondOk(w, r)
	}
}

type postTagsRequest struct {
	Tags []*models.Tag `json:"tags"`
}
