package common

import "testing"

func TestPickLimit(t *testing.T) {
	cases := []struct {
		name     string
		scene    int
		global   int
		fallback int
		want     int
	}{
		{"scene override", 5, 10, 20, 5},
		{"scene unlimited", -1, 10, 20, -1},
		{"global value", 0, 10, 20, 10},
		{"global unlimited", 0, -1, 20, -1},
		{"fallback", 0, 0, 20, 20},
	}
	for _, tc := range cases {
		if got := pickLimit(tc.scene, tc.global, tc.fallback); got != tc.want {
			t.Fatalf("%s: got %d, want %d", tc.name, got, tc.want)
		}
	}
}
