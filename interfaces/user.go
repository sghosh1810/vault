package interfaces

type UserCreatePayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserUpdatePayload struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type UserDeletePayload struct {
	ID int `json:"id"`
}

type UserSigninPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
