package main

import "testing"

func TestCleanInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "'  hello  world  ' should be ['hello','world']",
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			name:     "'Charmander Bulbasaur PIKACHU' should be ['charmander','bulbasaur','pikachu']",
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := cleanInput(tt.input)
			if len(actual) != len(tt.expected) {
				t.Errorf("lengths of slices do not match: expected %v != actual %v",
					tt.expected, actual,
				)
			}
			for i := range actual {
				word := actual[i]
				expectedWord := tt.expected[i]
				if word != expectedWord {
					t.Fatalf("expected word does not match - '%s' != '%s'", expectedWord, word)
				}
			}
		})
	}
}
