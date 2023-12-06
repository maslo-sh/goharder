package detection

type Detector interface {
	DetectMaliciousContent([]byte) int
}

type SqlDetector struct {
}

type LdapDetector struct {
}

func (d *SqlDetector) DetectMaliciousContent(payload []byte) int {
	return CANNOT_DEFINE
}

func (d *LdapDetector) DetectMaliciousContent(payload []byte) int {
	return CANNOT_DEFINE
}
