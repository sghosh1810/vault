package interfaces

type ProjectsCreatePayload struct {
	ProjectDisplayName string `json:"name"`
}

type ProjectsUpdatePayload struct {
	ProjectID          int    `json:"id"`
	ProjectDisplayName string `json:"name"`
}

type ProjectsDeletePayload struct {
	ProjectID int `json:"id"`
}
