package models

type User struct {
	ID       string `json:"id" dynamo:",hash"`
	Username string `json:"username" dynamo:""`
	Password string `json:"-" dynamo:""`
}
