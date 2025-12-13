package interfaces

type UserPayload struct {
	Uid            int64  `json:"uid"`
	SessionID      string `json:"session_id"`
	FirstName      string `json:"FirstName"`
	LastName       string `json:"LastName"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profilePicture"`
	Password       string `json:"Password"`
}

type UserSignupPayload struct {
	//Id             int64 `json:"id"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profilePicture"`
	Password       string `json:"password"`
}

type UserSignInPayload struct {
	Email    string `json:"email"`
	Password string `json:"Password"`
}
type UserCredentials struct {
	ID           int64
	PasswordHash string
}
