package hash

import (
	"testing"
)

func TestHash(t *testing.T) {
	var tests = []struct {
		name string
		v   string 
		h1 uint64
		h2 uint64
	}{
		{"Hash of short string", "hello world", 0xc05292b747fc78c0, 0x85bdab5e19e59315},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h1, h2, err := HashBytes([]byte(tt.v), 42)
			if err != nil {
				t.Errorf("Got error: %v", err)
			}
			if tt.h1 != h1 || tt.h2 != h2 {
				t.Errorf("Got (%x, %x) wanted (%x, %x)", h1, h2, tt.h1, tt.h2)
			}
		})
	}

}
