package sqlengine

const (
	ProjectInsertQuery = `
		INSERT INTO project (name,uid) VALUES (?,?)
	`
	ProjectUpdateQuery = `
		UPDATE project SET name = ? WHERE id = ?
	`
	ProjectGetQuery = `
		SELECT * from project WHERE id = ?
	`
	DeleteProjectFromMapTable = `
		DELETE FROM secret_project_map WHERE project_id = ?
	`
	DeleteProjectFromProjectTable = `
		DELETE FROM project WHERE id = ?
	`
	InsertUserProjectAccessQuery = `
		INSERT INTO user_project_access (
			user_id,
			project_id,
			has_write_access,
			has_share_access,
			has_delete_access
		) VALUES (?, ?, ?, ?, ?);
	`
	ListAllProjectByUser = `
		SELECT p.id,p.name
		FROM project p
		JOIN user_project_access upa
		ON p.id = upa.project_id
		WHERE upa.user_id = ?
	`
)

const (
	SecretInsertQuery = `
		INSERT INTO secret (name) VALUES (?);
	`
	SecretVersionInsertQuery = `
		INSERT INTO secret_version (secret_id, version, value, environment)
		VALUES (?, ?, ?, ?);
	`
	SecretGetLatestVersionQuery = `
		SELECT IFNULL(MAX(version), 0)
		FROM secret_version
		WHERE secret_id = ?;
	`
	SecretGetLatestValueQuery = `
		SELECT s.id, s.name, sv.version, sv.value, sv.environment
		FROM secret s
		JOIN secret_version sv ON s.id = sv.secret_id
		WHERE s.id = ?
		ORDER BY sv.version DESC
		LIMIT 1;
	`
	DeleteSecretFromMapTable = `
		DELETE FROM secret_project_map WHERE secret_id = ?
	`
	DeleteSecretFromSecretVersionTable = `
		DELETE FROM secret_version WHERE secret_id = ?
	`
	DeleteSecretFromSecretTable = `
		DELETE FROM secret WHERE id = ?
	`
	InsertUserSecretAccessQuery = `
		INSERT INTO user_secret_access (
			user_id,
			secret_id,
			has_write_access,
			has_share_access,
			has_delete_access
		) VALUES (?, ?, ?, ?, ?);
	`
	ListAllSecretByUser = `
		SELECT s.id,s.name
		FROM secret s
		JOIN user_secret_access usa
		ON s.id = usa.secret_id
		WHERE usa.user_id = ?
	`
)

const (
	GetAllSecretByProjectUidQuery = `
		SELECT s.id, s.name, sv.version, sv.value, sv.environment
		FROM secret s
		JOIN secret_version sv on sv.secret_id = s.id
		JOIN secret_project_map sp on sp.secret_id = s.id
		JOIN project p on p.id = sp.project_id
		WHERE p.uid = ? AND sv.environment = ?
		ORDER BY sv.id DESC 
	`
)

const (
	MapSecretToProject = `
		INSERT INTO secret_project_map (project_id,secret_id) VALUES (?, ?)
	`
)

const (
	CheckUserProjectAccessQuery = `
		SELECT has_read_access, has_write_access, has_share_access, has_delete_access
		FROM user_project_access
		WHERE user_id = ? AND project_id = ?;
	`
	CheckUserSecretAccessQuery = `
		SELECT has_read_access, has_write_access, has_share_access, has_delete_access
		FROM user_secret_access
		WHERE user_id = ? AND secret_id = ?;
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
		SELECT id, refresh_token_hash, refresh_token_expires_at
		FROM user_sessions
		WHERE user_id = ? AND session_id = ?;
	`
	UpdateUserSessionIDForUserQuery = `
		UPDATE user_sessions SET session_id = ?, refresh_token_hash = ?, refresh_token_expires_at = ?, updated_at = ? WHERE id = ?;
	`
	DeleteExpiredSessionsQuery = `DELETE FROM user_sessions WHERE refresh_token_expires_at < CURRENT_TIMESTAMP`
)
