package hash

import (
	"testing"
)

func TestConstants(t *testing.T) {
	var tests = []struct {
		name string
		v    uint64
		want uint64
	}{
		{"Constant c1 should be 0x87c37b91114253d5 for compatibility with Apache", c1, 0x87c37b91114253d5},
		{"Constant c2 should be 0x4cf5ad432745937f for compatibility with Apache", c2, 0x4cf5ad432745937f},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.v != tt.want {
				t.Errorf("Got %x, wanted %x", tt.v, tt.want)
			}
		})
	}

}
