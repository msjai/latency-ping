package pingpong

import (
	"net/http"

	"github.com/msjai/latency-ping/internal/config"
	"github.com/msjai/latency-ping/internal/entity"
)

type LatencyPingPongI interface {
	GetLatency(siteName string) (*entity.WebSite, error)
}

// LatencyPingPong -.
type LatencyPingPong struct {
	cfg *config.Config
}

// New -.
func New(config *config.Config) (*LatencyPingPong, error) {
	return &LatencyPingPong{cfg: config}, nil
}

func (pingPong *LatencyPingPong) GetLatency(siteName string) (*entity.WebSite, error) {
	addr := "https://" + siteName
	result := &entity.WebSite{Name: siteName, Address: addr}

	tp := newTransport()
	client := &http.Client{Transport: tp}

	resp, err := client.Get(addr)
	if err != nil {
		pingPong.cfg.L.Errorf("error get latency for %v: %v", addr, err)
		return result, err
	}
	defer resp.Body.Close()

	latency := tp.ConnDuration()
	result.Latency = &latency

	return result, nil
}
