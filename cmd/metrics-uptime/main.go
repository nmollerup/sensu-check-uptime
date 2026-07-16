package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	Scheme string
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "metrics-uptime",
			Short:    "Outputs system uptime and idle time metrics in Graphite plaintext format",
			Keyspace: "sensu.io/plugins/metrics-uptime/config",
		},
	}

	options = []sensu.ConfigOption{
		&sensu.PluginConfigOption[string]{
			Path:      "scheme",
			Argument:  "scheme",
			Shorthand: "s",
			Default:   defaultScheme(),
			Usage:     "Metric naming scheme, text to prepend to metric",
			Value:     &plugin.Scheme,
		},
	}
)

func defaultScheme() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	return hostname + ".uptime"
}

func main() {
	check := sensu.NewCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(_ *corev2.Event) (int, error) {
	return sensu.CheckStateOK, nil
}

// parseProcUptime parses the contents of /proc/uptime, which contains two
// space separated values: seconds since boot, and seconds spent idle.
func parseProcUptime(data []byte) (uptimeSec, idleSec float64, err error) {
	fields := strings.Fields(string(data))
	if len(fields) < 2 {
		return 0, 0, fmt.Errorf("unexpected /proc/uptime format: %q", string(data))
	}
	uptimeSec, err = strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("cannot parse uptime: %v", err)
	}
	idleSec, err = strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("cannot parse idle time: %v", err)
	}
	return uptimeSec, idleSec, nil
}

// formatMetrics renders the uptime and idle time metrics in the Graphite
// plaintext protocol: "path value timestamp".
func formatMetrics(scheme string, uptimeSec, idleSec float64, now time.Time) string {
	ts := now.Unix()
	return fmt.Sprintf("%s.uptime %d %d\n%s.idletime %d %d\n",
		scheme, int64(uptimeSec), ts, scheme, int64(idleSec), ts)
}

func executeCheck(_ *corev2.Event) (int, error) {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return sensu.CheckStateUnknown, fmt.Errorf("unable to read /proc/uptime: %v", err)
	}

	uptimeSec, idleSec, err := parseProcUptime(data)
	if err != nil {
		return sensu.CheckStateUnknown, err
	}

	fmt.Print(formatMetrics(plugin.Scheme, uptimeSec, idleSec, time.Now()))
	return sensu.CheckStateOK, nil
}
