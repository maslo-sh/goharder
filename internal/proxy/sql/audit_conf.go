package sql

type AuditConfiguration struct {
	ConfigFilePath      string
	CheckAuditLogs      bool
	CheckRemoteAccess   bool
	CheckAuditExtension bool
}
