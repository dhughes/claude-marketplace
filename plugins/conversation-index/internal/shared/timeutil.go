package shared

import (
	"fmt"
	"time"
)

// RFC3339Millis is the timestamp format used in Claude Code JSONL files
const RFC3339Millis = "2006-01-02T15:04:05.999Z07:00"

// ParseTimestamp parses various ISO 8601 timestamp formats
func ParseTimestamp(s string) (time.Time, error) {
	// Try RFC3339 with milliseconds
	if t, err := time.Parse(RFC3339Millis, s); err == nil {
		return t, nil
	}

	// Try RFC3339
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}

	// Try without timezone
	if t, err := time.Parse("2006-01-02T15:04:05", s); err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("unable to parse timestamp: %s", s)
}

// FormatTimestamp formats a time.Time to RFC3339 with milliseconds
func FormatTimestamp(t time.Time) string {
	return t.Format(RFC3339Millis)
}

// TruncateString safely truncates a UTF-8 string to maxLen characters (not bytes)
func TruncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-3]) + "..."
}
