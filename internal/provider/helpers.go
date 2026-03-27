package provider

import "strings"

func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	return strings.HasPrefix(err.Error(), "API result failed: 404")
}

func diffStringSets(oldSet, newSet []string) (toAdd, toRemove []string) {
	oldMap := make(map[string]bool, len(oldSet))
	for _, s := range oldSet {
		oldMap[s] = true
	}
	newMap := make(map[string]bool, len(newSet))
	for _, s := range newSet {
		newMap[s] = true
	}
	for _, s := range newSet {
		if !oldMap[s] {
			toAdd = append(toAdd, s)
		}
	}
	for _, s := range oldSet {
		if !newMap[s] {
			toRemove = append(toRemove, s)
		}
	}
	return
}
