package sql

const (
	SuperusersQuery      = "SELECT usename FROM pg_user WHERE usesuper = true;"
	PgauditQuery         = "SELECT * FROM pg_extension WHERE extname = 'pgaudit';"
	PgauditLoggingQuery  = "SELECT name, setting FROM pg_settings WHERE name = 'pgaudit.log';"
	ListenAddressesQuery = "SELECT setting FROM pg_settings WHERE name = 'listen_addresses';"

	AuthMethodQuery = "SELECT setting FROM pg_settings WHERE name = 'password_encryption';"
)
