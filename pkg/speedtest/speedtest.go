package speedtest

import (
	"fmt"

	st "github.com/showwin/speedtest-go/speedtest"
	"github.com/willfantom/sti/pkg/tester"
)

type Speedtest struct {
	serverID string
}

func New(serverID string) tester.Test {
	return &Speedtest{
		serverID: serverID,
	}
}

func (t *Speedtest) Name() string {
	return "speedtest"
}

func (t *Speedtest) Config() map[string]any {
	return map[string]any{
		"server_id": t.serverID,
	}
}

func (t *Speedtest) RunTest() (*tester.Result, error) {
	server, err := st.FetchServerByID(t.serverID)
	if err != nil {
		return nil, fmt.Errorf("speedtest server lookup failed: %w", err)
	}
	userInfo, err := st.FetchUserInfo()
	if err != nil {
		return nil, fmt.Errorf("speedtest user info fetch failed: %w", err)
	}
	if err := server.PingTest(nil); err != nil {
		return nil, fmt.Errorf("failed to run ping test: %w", err)
	}
	if err := server.DownloadTest(); err != nil {
		return nil, fmt.Errorf("failed to run download test: %w", err)
	}
	if err := server.UploadTest(); err != nil {
		return nil, fmt.Errorf("failed to run upload test: %w", err)
	}
	return &tester.Result{
		Labels: map[string]string{
			"server_id":      t.serverID,
			"server_country": server.Country,
			"user_ip":        userInfo.IP,
			"user_isp":       userInfo.Isp,
		},
		Data: map[string]any{
			"latency_ms":     server.Latency.Milliseconds(),
			"jitter_ms":      server.Jitter.Milliseconds(),
			"download_speed": server.DLSpeed.Mbps(),
			"upload_speed":   server.ULSpeed.Mbps(),
		},
	}, nil
}
