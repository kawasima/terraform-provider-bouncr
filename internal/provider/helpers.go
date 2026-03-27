package provider

import "strings"

func isNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), "404")
}
