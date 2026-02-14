package interfaces

type SecretCreatePayload struct {
	SecretName        string `json:"name"`
	SecretValue       string `json:"value"`
	SecretEnvironment string `json:"environment"`
	WorkspaceID       int    `json:"workspace_id"`
}

type SecretUpdatePayload struct {
	SecretID    int    `json:"id"`
	SecretValue string `json:"value"`
	WorkspaceID int    `json:"workspace_id"`
}

type SecretDeletePayload struct {
	SecretID    int `json:"id"`
	WorkspaceID int `json:"workspace_id"`
}

type SecretResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Version     int    `json:"version"`
	Value       string `json:"value"`
	Environment string `json:"environment"`
}
