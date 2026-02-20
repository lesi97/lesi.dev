package store

import "testing"

func TestNormaliseGameName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "removes trademark symbol",
			input:    "HELLDIVERS™ 2",
			expected: "helldivers 2",
		},
		{
			name:     "normalises punctuation and spacing",
			input:    "  Tom Clancy's  Rainbow: Six Siege ",
			expected: "tom clancy s rainbow six siege",
		},
		{
			name:     "removes registered symbol",
			input:    "Tom Clancy's Rainbow Six® Siege",
			expected: "tom clancy s rainbow six siege",
		},
		{
			name:     "empty when only symbols",
			input:    " ™®::: ",
			expected: "",
		},
	}

	for _, test := range tests {
		result := NormaliseGameName(test.input)
		if result != test.expected {
			t.Fatalf("%s: expected %q got %q", test.name, test.expected, result)
		}
	}
}
