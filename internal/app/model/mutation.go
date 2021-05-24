package model

type Mutation struct {
	Model
	DestinationID int     `json:"destination_id"`
	OriginID      int     `json:"origin_id"`
	Amount        float64 `json:"amount"`
	Status        int     `json:"status"`
	Pocket        int     `json:"pocket"`
	Channel       int     `json:"channel"`
}
