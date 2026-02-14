package interfaces

type MapOperationPayload struct {
	ProjectID   int `json:"project_id"`
	SecretID    int `json:"secret_id"`
	WorkspaceID int `json:"workspace_id"`
}
