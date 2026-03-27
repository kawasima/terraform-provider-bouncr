package provider

import (
	"fmt"
	"testing"
)

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"404 error", fmt.Errorf("API result failed: 404 Not Found: NotFound"), true},
		{"409 error", fmt.Errorf("API result failed: 409 Conflict: {\"message\":\"Conflict.\"}"), false},
		{"500 error", fmt.Errorf("API result failed: 500 Internal Server Error"), false},
		{"random error", fmt.Errorf("connection refused"), false},
		{"contains 404 but not prefix", fmt.Errorf("something 404 happened"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isNotFound(tt.err)
			if got != tt.expected {
				t.Errorf("isNotFound(%v) = %v, want %v", tt.err, got, tt.expected)
			}
		})
	}
}

func TestDiffStringSets(t *testing.T) {
	tests := []struct {
		name       string
		oldSet     []string
		newSet     []string
		wantAdd    []string
		wantRemove []string
	}{
		{
			name:       "no changes",
			oldSet:     []string{"a", "b"},
			newSet:     []string{"a", "b"},
			wantAdd:    nil,
			wantRemove: nil,
		},
		{
			name:       "add only",
			oldSet:     []string{"a"},
			newSet:     []string{"a", "b", "c"},
			wantAdd:    []string{"b", "c"},
			wantRemove: nil,
		},
		{
			name:       "remove only",
			oldSet:     []string{"a", "b", "c"},
			newSet:     []string{"a"},
			wantAdd:    nil,
			wantRemove: []string{"b", "c"},
		},
		{
			name:       "add and remove",
			oldSet:     []string{"a", "b"},
			newSet:     []string{"b", "c"},
			wantAdd:    []string{"c"},
			wantRemove: []string{"a"},
		},
		{
			name:       "both empty",
			oldSet:     []string{},
			newSet:     []string{},
			wantAdd:    nil,
			wantRemove: nil,
		},
		{
			name:       "old empty",
			oldSet:     []string{},
			newSet:     []string{"a"},
			wantAdd:    []string{"a"},
			wantRemove: nil,
		},
		{
			name:       "new empty",
			oldSet:     []string{"a"},
			newSet:     []string{},
			wantAdd:    nil,
			wantRemove: []string{"a"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAdd, gotRemove := diffStringSets(tt.oldSet, tt.newSet)
			if !sliceEqual(gotAdd, tt.wantAdd) {
				t.Errorf("toAdd = %v, want %v", gotAdd, tt.wantAdd)
			}
			if !sliceEqual(gotRemove, tt.wantRemove) {
				t.Errorf("toRemove = %v, want %v", gotRemove, tt.wantRemove)
			}
		})
	}
}

func sliceEqual(a, b []string) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	if len(a) != len(b) {
		return false
	}
	m := make(map[string]int)
	for _, s := range a {
		m[s]++
	}
	for _, s := range b {
		m[s]--
		if m[s] < 0 {
			return false
		}
	}
	return true
}
