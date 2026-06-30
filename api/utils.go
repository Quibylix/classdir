package main

import "regexp"

var uuidV7Regex = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func isValidUUIDv7(s string) bool {
	return uuidV7Regex.MatchString(s)
}
