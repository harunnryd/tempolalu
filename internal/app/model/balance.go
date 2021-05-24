package model

type Balance struct {
	Model
	UserID  int `json:"user_id"`
	Nominal int `json:"nominal"`
}
