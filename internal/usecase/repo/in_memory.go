package repo

import (
	"sync"

	"github.com/msjai/latency-ping/internal/config"
	"github.com/msjai/latency-ping/internal/entity"
	"github.com/msjai/latency-ping/internal/lib"
)

// LatencyRepoI - общий интерфейс репозитория
// В каждом пакете дополнительно объявлены минималистичные интерфейсы "провайдеры"
// по месту использования
type LatencyRepoI interface {
	GetMinLatencyUnit() *entity.WebSite
	GetMaxLatencyUnit() *entity.WebSite
	GetLatencyUnit(name string) *entity.WebSite
	UpdateLatencyUnits(webSite *entity.WebSite) error
	// .... more methods
}

// LatencyRepo .-
type LatencyRepo struct {
	mx         sync.RWMutex
	dB         map[string]*entity.WebSite
	minLatency *entity.WebSite
	maxLatency *entity.WebSite
}

// New -.
func New(cfg *config.Config) (*LatencyRepo, error) {
	cfg.L.Infow("Initializing Latency repository in memory...")

	db := make(map[string]*entity.WebSite, 50)
	repo := &LatencyRepo{dB: db}

	return repo, nil
}

func (repo *LatencyRepo) GetMinLatencyUnit() *entity.WebSite {
	repo.mx.RLock()
	defer repo.mx.RUnlock()

	return repo.minLatency
}

func (repo *LatencyRepo) GetMaxLatencyUnit() *entity.WebSite {
	repo.mx.RLock()
	defer repo.mx.RUnlock()

	return repo.maxLatency
}

func (repo *LatencyRepo) GetLatencyUnit(name string) (*entity.WebSite, bool) {
	repo.mx.RLock()
	defer repo.mx.RUnlock()
	val, ok := repo.dB[name]

	return val, ok
}

func (repo *LatencyRepo) UpdateLatencyUnits(webSite *entity.WebSite) error {
	repo.mx.Lock()
	defer repo.mx.Unlock()

	repo.dB[webSite.Name] = webSite

	if webSite.Latency != nil {
		if repo.minLatency == nil || lib.GetDurationFromPtr(webSite.Latency) < lib.GetDurationFromPtr(repo.minLatency.Latency) {
			repo.minLatency = webSite

		}
	}

	if webSite.Latency != nil {
		if repo.maxLatency == nil || lib.GetDurationFromPtr(webSite.Latency) > lib.GetDurationFromPtr(repo.maxLatency.Latency) {
			repo.maxLatency = webSite

		}
	}

	return nil
}
