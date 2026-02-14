package interfaces

type WorkspaceCreatePayload struct {
	WorkspaceDisplayName string `json:"name"`
	WorkspaceDescription string `json:"description"`
}

type WorkspaceUpdatePayload struct {
	WorkspaceID          int    `json:"id"`
	WorkspaceDisplayName string `json:"name"`
	WorkspaceDescription string `json:"description"`
}

type WorkspaceDeletePayload struct {
	WorkspaceID int `json:"id"`
}
