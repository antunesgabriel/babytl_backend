package auth

type AuthWithEmailDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
