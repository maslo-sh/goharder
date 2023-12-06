package sql

type AuditData struct {
	NumberOfSuperusers      int
	IsAuditExtensionEnabled bool
	IsAuditLoggingEnabled   bool
	DatabaseHosts           []string
}
