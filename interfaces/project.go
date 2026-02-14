package interfaces

type ProjectsCreatePayload struct {
	ProjectDisplayName string `json:"name"`
	WorkspaceID        int    `json:"workspace_id"`
}

type ProjectsUpdatePayload struct {
	ProjectID          int    `json:"id"`
	ProjectDisplayName string `json:"name"`
	WorkspaceID        int    `json:"workspace_id"`
}

type ProjectsDeletePayload struct {
	ProjectID   int `json:"id"`
	WorkspaceID int `json:"workspace_id"`
}
