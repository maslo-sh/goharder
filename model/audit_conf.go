package model

type AuditConfiguration struct {
	CheckAuditLogs            bool
	CheckRemoteAccess         bool
	CheckAuditExtension       bool
	CheckSuperusers           bool
	CheckAuthenticationMethod bool
}
