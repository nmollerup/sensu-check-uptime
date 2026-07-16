package main

import (
	"fmt"
	"time"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	Warn        int
	GreaterThan bool
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "check-uptime",
			Short:    "Checks system uptime and warns if the system has recently rebooted",
			Keyspace: "sensu.io/plugins/check-uptime/config",
		},
	}

	options = []sensu.ConfigOption{
		&sensu.PluginConfigOption[int]{
			Path:      "warn",
			Argument:  "warn",
			Shorthand: "w",
			Default:   180,
			Usage:     "Warn threshold in seconds",
			Value:     &plugin.Warn,
		},
		&sensu.PluginConfigOption[bool]{
			Path:      "greater-than",
			Argument:  "greater-than",
			Shorthand: "g",
			Default:   false,
			Usage:     "Compare uptime > threshold. Default behavior is uptime < threshold",
			Value:     &plugin.GreaterThan,
		},
	}
)

func main() {
	check := sensu.NewCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(_ *corev2.Event) (int, error) {
	if plugin.Warn <= 0 {
		return sensu.CheckStateWarning, fmt.Errorf("--warn must be greater than 0")
	}
	return sensu.CheckStateOK, nil
}

// evaluateUptime contains the comparison logic, kept separate from
// getUptimeSeconds so it can be tested without depending on the host's
// actual uptime.
func evaluateUptime(uptimeSec float64, warn int, greaterThan bool) (int, string) {
	bootTime := time.Now().Add(-time.Duration(uptimeSec) * time.Second)

	if greaterThan && uptimeSec > float64(warn) {
		return sensu.CheckStateWarning, fmt.Sprintf("warning: system boot detected (%d seconds up), compared using '>'", int(uptimeSec))
	}
	if !greaterThan && uptimeSec < float64(warn) {
		return sensu.CheckStateWarning, fmt.Sprintf("warning: system boot detected (%d seconds up), compared using '<'", int(uptimeSec))
	}

	return sensu.CheckStateOK, fmt.Sprintf("ok: system booted at %v", bootTime.Format(time.RFC1123))
}

func executeCheck(_ *corev2.Event) (int, error) {
	uptimeSec, err := getUptimeSeconds()
	if err != nil {
		return sensu.CheckStateUnknown, fmt.Errorf("unable to determine system uptime: %v", err)
	}

	status, message := evaluateUptime(uptimeSec, plugin.Warn, plugin.GreaterThan)
	fmt.Println(message)
	return status, nil
}
