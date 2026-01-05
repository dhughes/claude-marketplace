package shared

import "strings"

// EncodeProjectPath converts a file path to an encoded directory name
// /Users/foo/bar -> -Users-foo-bar
func EncodeProjectPath(path string) string {
	if strings.HasPrefix(path, "/") {
		return "-" + strings.ReplaceAll(path[1:], "/", "-")
	}
	return strings.ReplaceAll(path, "/", "-")
}

// DecodeProjectPath converts an encoded directory name back to a path
// -Users-foo-bar -> /Users/foo/bar
// WARNING: This is lossy if directory names contain dashes
func DecodeProjectPath(encoded string) string {
	if !strings.HasPrefix(encoded, "-") {
		return encoded
	}
	return "/" + strings.ReplaceAll(encoded[1:], "-", "/")
}
