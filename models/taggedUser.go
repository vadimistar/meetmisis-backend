package models

type TaggedUser struct {
	UserID string   `json:"userID" dynamo:",hash"`
	Tags   []string `json:"tags" dynamo:""`
}
