// Package rediskeys provides a global Redis key prefix for shared instances (e.g. Upstash).
package rediskeys

import "strings"

var prefix string

// Init sets the global key prefix (normalized to end with ":"). Call once at startup.
func Init(raw string) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		prefix = ""
		return
	}
	if !strings.HasSuffix(raw, ":") {
		raw += ":"
	}
	prefix = raw
}

// Prefix returns the normalized global prefix (e.g. "vet-clinic-api:").
func Prefix() string {
	return prefix
}

// Key prepends the global prefix to a logical key suffix.
func Key(suffix string) string {
	if prefix == "" {
		return suffix
	}
	if strings.HasPrefix(suffix, prefix) {
		return suffix
	}
	return prefix + suffix
}
