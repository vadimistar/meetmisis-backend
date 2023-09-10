package models

type Tag struct {
	Id   string `json:"id,omitempty" dynamo:",hash"`
	Name string `json:"name,omitempty" dynamo:""`
}
