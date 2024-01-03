package relational

type AuditData struct {
	Superusers              []string
	IsAuditExtensionEnabled bool
	IsAuditLoggingEnabled   bool
	AuthenticationMethod    string
	DatabaseHosts           []string
}
