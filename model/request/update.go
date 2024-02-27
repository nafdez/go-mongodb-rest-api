package request

type Update struct {
	Points int32  `json:"points"`
	Token  string `json:"token"`
}
