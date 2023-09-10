package models

type Partner struct {
	UserID   string   `json:"userID" dynamo:",hash"`
	Partners []string `json:"partners" dynamo:""`
}
