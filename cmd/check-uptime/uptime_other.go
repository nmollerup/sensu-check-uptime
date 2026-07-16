//go:build !linux && !windows

package main

import (
	"fmt"
	"runtime"
)

// getUptimeSeconds is unsupported on platforms other than linux and windows.
func getUptimeSeconds() (float64, error) {
	return 0, fmt.Errorf("platform %s is not supported, please open an issue", runtime.GOOS)
}
