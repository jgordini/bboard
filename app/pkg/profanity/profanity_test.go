package profanity

import (
	"testing"
)

func TestContainsProfanity(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		expect bool
	}{
		{"empty", "", false},
		{"clean", "hello world", false},
		{"clean sentence", "This is a normal title for a post.", false},
		{"classic not matched", "classic design", false},
		{"assassin not matched", "assassin", false},
		{"dirty word", "this is shit", true},
		{"dirty word caps", "FUCK this", true},
		{"dirty word mid", "something fuck something", true},
		{"slur", "that nigger", true},
		{"slur variant", "nigga please", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := ContainsProfanity(tt.text)
			if got != tt.expect {
				t.Errorf("ContainsProfanity(%q) = %v, want %v", tt.text, got, tt.expect)
			}
		})
	}
}
