package models

type Verification struct {
	UserID string `json:"userID" dymano:",hash" db:"userID"`
	Token  string `json:"-" dymano:"" db:"token"`
}
