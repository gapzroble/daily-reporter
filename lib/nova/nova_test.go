package nova

import "testing"

// TestSetNovaHours test
func TestSetNovaHours(t *testing.T) {
	tests := []struct {
		input  float64
		output string
	}{
		{8.5, "8,5"},
		{8.0, "8"},
		{7, "7"},
		{6.25, "6,25"},
	}

	for _, test := range tests {
		setNovaHours(test.input)
		if novaHours != test.output {
			t.Errorf("Expecting novaHours %.2f = '%s', got '%s'", test.input, test.output, novaHours)
		}
	}
}
