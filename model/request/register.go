package request

type Register struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email,"`
	Name     string `json:"name,omitempty"`
	Role     string `json:"role,omitempty"`
}
