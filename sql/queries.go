package sqlengine

const (
	ProjectInsertQuery = `
		INSERT INTO project (name,uid,workspace_id) VALUES (?,?,?);
	`
	ProjectUpdateQuery = `
		UPDATE project SET name = ? WHERE id = ? AND workspace_id = ?
	`
	ProjectGetQuery = `
		SELECT p.id,p.name,p.uid from project p WHERE id = ? AND workspace_id = ?
	`
	DeleteProjectFromMapTable = `
		DELETE FROM secret_project_map WHERE project_id = ? 
	`
	DeleteProjectFromProjectTable = `
		DELETE FROM project WHERE id = ? AND workspace_id = ?
	`
	InsertWorkspaceProjectAccessQuery = `
		INSERT INTO workspace_project_access (
			workspace_id,
			project_id,
			has_write_access,
			has_share_access,
			has_delete_access
		) VALUES (?, ?, ?, ?, ?);
	`
	ListAllProjectByWorkspace = `
		SELECT p.id,p.name,p.uid
		FROM project p
		JOIN workspace_project_access wpa ON p.id = wpa.project_id
		JOIN user_workspace_access uwa ON wpa.workspace_id = uwa.id
		WHERE uwa.user_id = ?
		AND uwa.workspace_id = ?
	`
)

const (
	SecretInsertQuery = `
		INSERT INTO secret (name, workspace_id) VALUES (?, ?);
	`
	SecretVersionInsertQuery = `
		INSERT INTO secret_version (secret_id, version, value, environment)
		VALUES (?, ?, ?, ?);
	`
	SecretGetLatestVersionQuery = `
		SELECT IFNULL(MAX(version), 0)
		FROM secret_version
		WHERE secret_id = ?
		AND environment = ?
	`
	SecretGetLatestValueQuery = `
		SELECT id, name, version, value, environment
		FROM (
			SELECT s.id,
				s.name,
				sv.version,
				sv.value,
				sv.environment,
				ROW_NUMBER() OVER (
					PARTITION BY sv.environment
					ORDER BY sv.version DESC
				) as rn
			FROM secret s
			JOIN secret_version sv ON s.id = sv.secret_id
			WHERE s.id = ?
		) t
		WHERE rn = 1;
	`
	DeleteSecretFromMapTable = `
		DELETE FROM secret_project_map WHERE secret_id = ?
	`
	DeleteSecretFromSecretVersionTable = `
		DELETE FROM secret_version WHERE secret_id = ?
	`
	DeleteSecretFromSecretTable = `
		DELETE FROM secret WHERE id = ? AND workspace_id = ?
	`
	InsertWorkspaceSecretAccessQuery = `
		INSERT INTO workspace_secret_access (
			workspace_id,
			secret_id,
			has_write_access,
			has_share_access,
			has_delete_access
		) VALUES (?, ?, ?, ?, ?);
	`
	ListAllSecretByWorkspace = `
		SELECT s.id,s.name
		FROM secret s
		JOIN workspace_secret_access wsa ON s.id = wsa.secret_id
		JOIN user_workspace_access uwa ON wsa.workspace_id = uwa.id
		WHERE uwa.user_id = ?
		AND uwa.workspace_id = ?
	`
)

const (
	GetAllSecretByProjectIdQuery = `
		SELECT s.id, s.name, sv.version, sv.environment
		FROM secret s
		JOIN secret_version sv on sv.secret_id = s.id
		JOIN secret_project_map sp on sp.secret_id = s.id
		JOIN project p on p.id = sp.project_id
		WHERE p.id = ?
		GROUP BY sv.environment
		ORDER BY sv.id DESC
	`
)

const (
	MapSecretToProject = `
		INSERT INTO secret_project_map (project_id,secret_id) VALUES (?, ?)
	`
)

const (
	CheckUserWorkspaceProjectAccessQuery = `
		SELECT uwa.has_read_access * wpa.has_read_access as has_read_access, uwa.has_write_access * wpa.has_write_access as has_write_access, uwa.has_share_access * wpa.has_share_access as has_share_access, uwa.has_delete_access * wpa.has_delete_access as has_delete_access
		FROM workspace_project_access wpa
		JOIN user_workspace_access uwa
		ON wpa.workspace_id = uwa.id
		WHERE uwa.user_id = ? 
		AND wpa.project_id = ? 
		AND wpa.workspace_id = ?;
	`
	CheckUserWorkspaceSecretAccessQuery = `
		SELECT uwa.has_read_access * wsa.has_read_access as has_read_access, uwa.has_write_access * wsa.has_write_access as has_write_access, uwa.has_share_access * wsa.has_share_access as has_share_access, uwa.has_delete_access * wsa.has_delete_access as has_delete_access
		FROM workspace_secret_access wsa
		JOIN user_workspace_access uwa
		ON wsa.workspace_id = uwa.id
		WHERE uwa.user_id = ?
		AND wsa.secret_id =  ?
		AND wsa.workspace_id =  ?;
	`
	CheckUserWorkspaceAccessQuery = `
		SELECT has_read_access, has_write_access, has_share_access, has_delete_access
		FROM user_workspace_access
		WHERE user_id = ? AND workspace_id = ?;
	`
)

const (
	InsertUserQuery = `
		INSERT INTO users (
			first_name, 
			last_name, 
			email, 
			password_hash, 
			profile_picture
		) VALUES (?, ?, ?, ?, ?);
	`
	CheckUserCreDentialsQuery = `
		SELECT id, password_hash
		FROM users
		WHERE email = ?;
	`
	InsertUserSessionsQuery = `
		INSERT INTO user_sessions (
			user_id, 
			session_id,
			refresh_token_hash,
			refresh_token_expires_at
		) VALUES (?, ?, ?, ?);
	`
	DeleteSessionIDFromUserSessionsTable = `
		DELETE FROM user_sessions
        WHERE user_id = ?
        AND session_id = ?;
	`
	SelectUserSessionQuery = `
		SELECT id, user_id, refresh_token_hash, refresh_token_expires_at
        FROM user_sessions
        WHERE session_id = ?
	`
	UpdateUserSessionIDForUserQuery = `
		UPDATE user_sessions SET session_id = ?, refresh_token_hash = ?, refresh_token_expires_at = ?, updated_at = ? WHERE id = ?;
	`
	DeleteExpiredSessionsQuery = `DELETE FROM user_sessions WHERE refresh_token_expires_at < CURRENT_TIMESTAMP`
)

const (
	WorkspaceInsertQuery = `
		INSERT INTO workspace (name, description) VALUES (?, ?)
	`
	InsertUserWorkspaceAccessQuery = `
		INSERT INTO user_workspace_access (
			user_id,
			workspace_id,
			has_read_access,
			has_write_access,
			has_share_access,
			has_delete_access
		) VALUES (?, ?, ?, ?, ?, ?);
	`
	WorkspaceUpdateQuery = `
		UPDATE workspace SET name = ?, description = ? WHERE id = ?
	`
	WorkspaceGetQuery = `
		SELECT id, name, description FROM workspace WHERE id = ?
	`

	DeleteWorkspaceFromWorkspaceTable = `
		DELETE FROM workspace WHERE id = ?
	`
	ListAllWorkspaceByUser = `
		SELECT w.id,w.name,w.description
		FROM workspace w
		JOIN user_workspace_access uwa
		ON w.id = uwa.workspace_id
		WHERE uwa.user_id = ?
	`
)
