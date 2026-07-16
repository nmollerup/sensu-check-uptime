//go:build windows

package main

import "syscall"

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	procGetTickCount64 = kernel32.NewProc("GetTickCount64")
)

// getUptimeSeconds reads system uptime via the GetTickCount64 Windows API.
func getUptimeSeconds() (float64, error) {
	r, _, _ := procGetTickCount64.Call()
	return float64(r) / 1000.0, nil
}
