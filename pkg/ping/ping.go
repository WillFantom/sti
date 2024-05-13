package ping

import (
	"fmt"

	pg "github.com/prometheus-community/pro-bing"
	"github.com/willfantom/influx-speedtest/pkg/tester"
)

// Ping represents a ping test.
type Ping struct {
	Target string // Target is the IP address or hostname to ping.
	Count  int    // Count is the number of pings to send.
}

// New creates a new Ping test instance.
func New(target string, count int) tester.Test {
	return &Ping{
		Target: target,
		Count:  count,
	}
}

// Name returns the name of the Ping test.
func (t *Ping) Name() string {
	return "ping"
}

// Config returns the configuration of the Ping test.
func (t *Ping) Config() map[string]interface{} {
	return map[string]interface{}{
		"target": t.Target,
		"count":  t.Count,
	}
}

// RunTest runs the Ping test and returns the result.
func (t *Ping) RunTest() (*tester.Result, error) {
	pinger, err := pg.NewPinger(t.Target)
	if err != nil {
		return nil, fmt.Errorf("failed to create pinger instance: %w", err)
	}
	pinger.Count = t.Count
	err = pinger.Run()
	if err != nil {
		return nil, fmt.Errorf("ping test exited with error: %w", err)
	}
	stats := pinger.Statistics()
	return &tester.Result{
		Labels: map[string]string{
			"target": t.Target,
		},
		Data: map[string]any{
			"packets_sent":     stats.PacketsSent,
			"packets_received": stats.PacketsRecv,
			"packet_loss":      stats.PacketLoss,
			"rtt_min_ms":       stats.MinRtt.Milliseconds(),
			"rtt_max_ms":       stats.MaxRtt.Milliseconds(),
			"rtt_avg_ms":       stats.AvgRtt.Milliseconds(),
			"rtt_mdev_ms":      stats.StdDevRtt.Milliseconds(),
		},
	}, nil
}
