package main

import (
	"strings"
	"testing"

	"github.com/sensu/sensu-plugin-sdk/sensu"
)

func TestCheckArgs(t *testing.T) {
	tests := []struct {
		name       string
		warn       int
		wantStatus int
		wantErr    bool
	}{
		{"zero warn", 0, sensu.CheckStateWarning, true},
		{"negative warn", -10, sensu.CheckStateWarning, true},
		{"valid warn", 180, sensu.CheckStateOK, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plugin = Config{Warn: tt.warn}

			status, err := checkArgs(nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if status != tt.wantStatus {
				t.Errorf("checkArgs() status = %v, want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestEvaluateUptime(t *testing.T) {
	tests := []struct {
		name        string
		uptimeSec   float64
		warn        int
		greaterThan bool
		wantStatus  int
		wantSub     string
	}{
		{"recent boot, less-than mode", 30, 180, false, sensu.CheckStateWarning, "boot detected"},
		{"long uptime, less-than mode", 3600, 180, false, sensu.CheckStateOK, "system booted at"},
		{"boundary uptime, less-than mode is ok", 180, 180, false, sensu.CheckStateOK, "system booted at"},
		{"long uptime, greater-than mode triggers warning", 3600, 180, true, sensu.CheckStateWarning, "boot detected"},
		{"short uptime, greater-than mode is ok", 30, 180, true, sensu.CheckStateOK, "system booted at"},
		{"boundary uptime, greater-than mode is ok", 180, 180, true, sensu.CheckStateOK, "system booted at"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, message := evaluateUptime(tt.uptimeSec, tt.warn, tt.greaterThan)
			if status != tt.wantStatus {
				t.Errorf("evaluateUptime() status = %v, want %v", status, tt.wantStatus)
			}
			if !strings.Contains(message, tt.wantSub) {
				t.Errorf("evaluateUptime() message = %q, want it to contain %q", message, tt.wantSub)
			}
		})
	}
}

func TestExecuteCheck(t *testing.T) {
	plugin = Config{Warn: 180, GreaterThan: false}

	status, err := executeCheck(nil)
	if err != nil {
		t.Fatalf("executeCheck() unexpected error: %v", err)
	}
	if status != sensu.CheckStateOK && status != sensu.CheckStateWarning {
		t.Errorf("executeCheck() status = %v, want OK or Warning", status)
	}
}

func TestGetUptimeSeconds(t *testing.T) {
	uptimeSec, err := getUptimeSeconds()
	if err != nil {
		t.Fatalf("getUptimeSeconds() unexpected error: %v", err)
	}
	if uptimeSec <= 0 {
		t.Errorf("getUptimeSeconds() = %v, want a positive value", uptimeSec)
	}
}
