package fs

// short truncates a session ID for log output. Session IDs are
// cryptographically random and long; the first 6 characters are
// sufficient for correlation in debug logs.
func short(id string) string {
	if len(id) > 6 {
		return id[:6]
	}
	return id
}
