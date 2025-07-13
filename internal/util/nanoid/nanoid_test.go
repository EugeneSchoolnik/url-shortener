package nanoid

import (
	"testing"
)

func TestID(t *testing.T) {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	tests := []struct {
		name     string
		alphabet string
		size     int
		wantLen  int
		wantErr  bool
	}{
		{
			name:     "1",
			alphabet: alphabet,
			size:     12,
			wantErr:  false,
		},
		{
			name:     "2",
			alphabet: alphabet,
			size:     8,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idGenerator := New(tt.alphabet, tt.size)

			got, err := idGenerator.ID()
			t.Log(got)

			if (err != nil) != tt.wantErr {
				t.Errorf("ID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.size {
				t.Errorf("ID() = %v, size %v", got, tt.size)
			}
		})
	}
}
