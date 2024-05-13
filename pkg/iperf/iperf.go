package iperf

import (
	"fmt"
	"net/netip"

	iperfcli "github.com/BGrewell/go-iperf"
	"github.com/willfantom/influx-speedtest/pkg/tester"
)

type Iperf struct {
	ServerIP   string
	ServerPort int
	Streams    int
	Seconds    int
	TCP        bool
}

func New(serverIP string, serverPort int, streams int, seconds int, tcp bool) tester.Test {
	return &Iperf{
		ServerIP:   serverIP,
		ServerPort: serverPort,
		Streams:    streams,
		Seconds:    seconds,
		TCP:        tcp,
	}
}

func (t *Iperf) Name() string {
	return "iperf"
}

func (t *Iperf) Config() map[string]any {
	return map[string]any{
		"server_ip":   t.ServerIP,
		"server_port": t.ServerPort,
		"streams":     t.Streams,
		"seconds":     t.Seconds,
		"tcp":         t.TCP,
	}
}

func (t *Iperf) RunTest() (*tester.Result, error) {
	if _, err := netip.ParseAddr(t.ServerIP); err != nil {
		return nil, fmt.Errorf("iperf server ip is invalid: %w", err)
	}
	c := iperfcli.NewClient(t.ServerIP)
	c.SetJSON(true)
	c.SetIncludeServer(true)
	c.SetStreams(t.Streams)
	c.SetTimeSec(t.Seconds)
	c.SetInterval(1)
	if t.TCP {
		c.SetProto(iperfcli.PROTO_TCP)
	} else {
		c.SetProto(iperfcli.PROTO_UDP)
	}
	c.SetPort(t.ServerPort)
	if err := c.Start(); err != nil {
		return nil, fmt.Errorf("iperf test failed to start: %w", err)
	}
	<-c.Done
	if c.Report().Error != "" || len(c.Report().Start.Connected) == 0 {
		return nil, fmt.Errorf("iperf test failed: %s", c.Report().Error)
	}

	return &tester.Result{
		Labels: map[string]string{
			"server_ip":   t.ServerIP,
			"server_port": fmt.Sprintf("%d", t.ServerPort),
			"streams":     fmt.Sprintf("%d", t.Streams),
			"seconds":     fmt.Sprintf("%d", t.Seconds),
			"tcp":         fmt.Sprintf("%t", t.TCP),
		},
		Data: map[string]any{
			"received_megabits_per_second": (c.Report().End.SumReceived.BitsPerSecond / 1000000),
		},
	}, nil
}
