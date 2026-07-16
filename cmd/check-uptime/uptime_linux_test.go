//go:build linux

package main

import "testing"

func TestParseProcUptime(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    float64
		wantErr bool
	}{
		{"typical", "12345.67 6789.01\n", 12345.67, false},
		{"single field", "12345.67\n", 12345.67, false},
		{"empty", "", 0, true},
		{"garbage", "not-a-number 123\n", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseProcUptime([]byte(tt.data))
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseProcUptime() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("parseProcUptime() = %v, want %v", got, tt.want)
			}
		})
	}
}
