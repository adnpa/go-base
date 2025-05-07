package tar

import (
	"testing"
)

func TestCompressPath(t *testing.T) {
	tests := []string{"./test"}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			if err := CompressPath(tt); err != nil {
				t.Error(err)
			}
		})
	}
}
