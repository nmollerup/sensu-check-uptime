package main

import (
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/sensu/sensu-plugin-sdk/sensu"
)

func TestCheckArgs(t *testing.T) {
	status, err := checkArgs(nil)
	if err != nil {
		t.Fatalf("checkArgs() unexpected error: %v", err)
	}
	if status != sensu.CheckStateOK {
		t.Errorf("checkArgs() status = %v, want OK", status)
	}
}

func TestDefaultScheme(t *testing.T) {
	scheme := defaultScheme()
	if !strings.HasSuffix(scheme, ".uptime") {
		t.Errorf("defaultScheme() = %q, want it to end with %q", scheme, ".uptime")
	}
}

func TestParseProcUptime(t *testing.T) {
	tests := []struct {
		name        string
		data        string
		wantUptime  float64
		wantIdle    float64
		wantErr     bool
		errContains string
	}{
		{"typical", "12345.67 6789.01\n", 12345.67, 6789.01, false, ""},
		{"empty", "", 0, 0, true, "unexpected"},
		{"single field", "12345.67\n", 0, 0, true, "unexpected"},
		{"non-numeric uptime", "abc 6789.01\n", 0, 0, true, "cannot parse uptime"},
		{"non-numeric idle", "12345.67 abc\n", 0, 0, true, "cannot parse idle time"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uptime, idle, err := parseProcUptime([]byte(tt.data))
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseProcUptime() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("parseProcUptime() error = %q, want it to contain %q", err.Error(), tt.errContains)
				}
				return
			}
			if uptime != tt.wantUptime {
				t.Errorf("parseProcUptime() uptime = %v, want %v", uptime, tt.wantUptime)
			}
			if idle != tt.wantIdle {
				t.Errorf("parseProcUptime() idle = %v, want %v", idle, tt.wantIdle)
			}
		})
	}
}

func TestFormatMetrics(t *testing.T) {
	now := time.Unix(1700000000, 0)
	got := formatMetrics("myhost.uptime", 12345, 6789, now)
	want := "myhost.uptime.uptime 12345 1700000000\nmyhost.uptime.idletime 6789 1700000000\n"
	if got != want {
		t.Errorf("formatMetrics() = %q, want %q", got, want)
	}
}

func TestExecuteCheck(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("executeCheck reads /proc/uptime, which only exists on linux")
	}
	if _, err := os.Stat("/proc/uptime"); err != nil {
		t.Skip("/proc/uptime not available")
	}

	plugin = Config{Scheme: "test.uptime"}
	status, err := executeCheck(nil)
	if err != nil {
		t.Fatalf("executeCheck() unexpected error: %v", err)
	}
	if status != sensu.CheckStateOK {
		t.Errorf("executeCheck() status = %v, want OK", status)
	}
}
