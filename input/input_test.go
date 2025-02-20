package input

import "testing"

func TestParseYesNoResponse(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		defaultValue bool
		want         bool
	}{
		{"Empty input with default true", "", true, true},
		{"Empty input with default false", "", false, false},
		{"Yes lowercase", "yes", false, true},
		{"Y lowercase", "y", false, true},
		{"Yes uppercase", "YES", false, true},
		{"Y uppercase", "Y", false, true},
		{"No lowercase", "no", true, false},
		{"N lowercase", "n", true, false},
		{"No uppercase", "NO", true, false},
		{"N uppercase", "N", true, false},
		{"Invalid input", "invalid", true, false},
		{"Whitespace with yes", "  yes  ", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseYesNoResponse(tt.input, tt.defaultValue)
			if got != tt.want {
				t.Errorf("ParseYesNoResponse(%q, %v) = %v, want %v",
					tt.input, tt.defaultValue, got, tt.want)
			}
		})
	}
}

func TestFormatYesNoPrompt(t *testing.T) {
	tests := []struct {
		name       string
		prompt     string
		defaultYes bool
		want       string
	}{
		{
			name:       "Default yes",
			prompt:     "Continue?",
			defaultYes: true,
			want:       "Continue? [Y/n]: ",
		},
		{
			name:       "Default no",
			prompt:     "Continue?",
			defaultYes: false,
			want:       "Continue? [y/N]: ",
		},
		{
			name:       "Empty prompt with default yes",
			prompt:     "",
			defaultYes: true,
			want:       " [Y/n]: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatYesNoPrompt(tt.prompt, tt.defaultYes)
			if got != tt.want {
				t.Errorf("FormatYesNoPrompt(%q, %v) = %q, want %q",
					tt.prompt, tt.defaultYes, got, tt.want)
			}
		})
	}
}
