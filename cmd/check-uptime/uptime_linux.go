//go:build linux

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// getUptimeSeconds reads system uptime from /proc/uptime.
func getUptimeSeconds() (float64, error) {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}
	return parseProcUptime(data)
}

func parseProcUptime(data []byte) (float64, error) {
	fields := strings.Fields(string(data))
	if len(fields) < 1 {
		return 0, fmt.Errorf("unexpected /proc/uptime format: %q", string(data))
	}
	return strconv.ParseFloat(fields[0], 64)
}
