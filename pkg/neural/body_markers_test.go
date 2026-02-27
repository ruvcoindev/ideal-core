package neural

import "testing"

func TestDetectBodyMarker(t *testing.T) {
	tests := []struct {
		text     string
		expected BodyMarker
		found    bool
	}{
		{"писю лечила", MarkerIntimacy, true},
		{"я прокакалась", MarkerCleansing, true},
		{"голова болит", MarkerPain, true},
		{"привет", "", false},
	}
	for _, tc := range tests {
		marker, found := DetectBodyMarker(tc.text)
		if found != tc.found || marker != tc.expected {
			t.Errorf("Text %q: expected (%s, %v), got (%s, %v)",
				tc.text, tc.expected, tc.found, marker, found)
		}
	}
}
