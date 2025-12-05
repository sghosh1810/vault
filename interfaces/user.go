package interfaces

type UserPayload struct {
	Uid            string `json:"uid"`
	PersonId       string `json:"personId"`
	ProfilePicture string `json:"profilePicture"`
	ExpiresAt      string `json:"expiresAt"`
	UserFullName   string `json:"userFullName"`
	UserName       string `json:"userName"`
}
