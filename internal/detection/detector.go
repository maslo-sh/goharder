package detection

import (
	"log"
	"strings"
)

type Detector interface {
	GetMaliciousQueries() []string
	DetectMaliciousContent([]byte) int
}

type SqlDetector struct {
}

type LdapDetector struct {
}

func (d SqlDetector) GetMaliciousQueries() []string {
	return []string{
		"union all select",
		"pg_read_file",
		"COPY file_store",
		"open(FD,\"$_[0] |\")",
		"whoami",
		"version()",
		"pg_sleep",
		"select current_user",
		"select session_user",
		"getpgusername()",
		"null--",
		"chr(",
		"ascii(",
	}
}

func (d SqlDetector) DetectMaliciousContent(payload []byte) int {
	query := strings.ToLower(string(payload))
	log.Printf("Packet with DQL type: %s", query)
	for _, v := range d.GetMaliciousQueries() {
		if strings.Contains(query, v) {
			return MALICIOUS
		}
	}

	return SAFE
}

func (d LdapDetector) DetectMaliciousContent(payload []byte) int {
	return CANNOT_DEFINE
}
