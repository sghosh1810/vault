package interfaces

type AccessControlPayload struct {
	HasReadAccess   bool `db:"has_read_access" json:"has_read_access"`
	HasWriteAccess  bool `db:"has_write_access" json:"has_write_access"`
	HasShareAccess  bool `db:"has_share_access" json:"has_share_access"`
	HasDeleteAccess bool `db:"has_delete_access" json:"has_delete_access"`
}
