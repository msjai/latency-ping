package usecase

import (
	"fmt"
	"time"

	"github.com/msjai/latency-ping/internal/config"
	"github.com/msjai/latency-ping/internal/entity"
)

// LatencyRepoProviderI - Интерфейс репозитория дополнительно объявляем по месту использования
// для уменьшения зависимостей и упрощения написания моков для юнит тестов.
// На уровне репозитория объявлен аналогичный интерфейс, но предполагается что на уровне репозитория
// интерфейс будет иметь больше методов чем здесь
type LatencyRepoProviderI interface {
	GetMinLatencyUnit() *entity.WebSite
	GetMaxLatencyUnit() *entity.WebSite
	GetLatencyUnit(name string) (*entity.WebSite, bool)
	UpdateLatencyUnits(webSite *entity.WebSite) error
}

// PingPongProviderI -
type PingPongProviderI interface {
	GetLatency(siteName string) (*entity.WebSite, error)
}

// LatencyUseCase -.
type LatencyUseCase struct {
	repo     LatencyRepoProviderI
	pingPong PingPongProviderI
	cfg      *config.Config
}

// New -.
func New(repo LatencyRepoProviderI, pingPong PingPongProviderI, cfg *config.Config) *LatencyUseCase {
	return &LatencyUseCase{
		repo:     repo,
		pingPong: pingPong,
		cfg:      cfg,
	}
}

// RefreshLatencyInfo - Здесь применяется паттерн многопоточности "WorkerPool"
func (latencyUseCase *LatencyUseCase) RefreshLatencyInfo() {
	workersCount := latencyUseCase.cfg.WorkerCount
	jobCh := make(chan *string) // канал содержит имена сайтов для получения статистики

	for i := 0; i < workersCount; i++ {
		go func() {
			for jobString := range jobCh {
				makeRefresh(jobString, latencyUseCase)
			}
		}()
	}

	for {
		for i := 0; i < len(latencyUseCase.cfg.ListWebSites); i++ {
			job := &latencyUseCase.cfg.ListWebSites[i]
			jobCh <- job
		}
		time.Sleep(60 * time.Second)
	}
}

func makeRefresh(siteName *string, latencyUseCase *LatencyUseCase) {
	webSite, _ := latencyUseCase.pingPong.GetLatency(*siteName)
	err := latencyUseCase.repo.UpdateLatencyUnits(webSite)
	if err != nil {
		latencyUseCase.cfg.L.Errorf("error updating latency for %v", webSite.Name)
	}

}

func (latencyUseCase *LatencyUseCase) GetMinLatencyLogic() (*entity.WebSite, error) {
	website := latencyUseCase.repo.GetMinLatencyUnit()
	if website == nil || website.Latency == nil {
		return nil, fmt.Errorf("no data")
	}

	return website, nil
}

func (latencyUseCase *LatencyUseCase) GetMaxLatencyLogic() (*entity.WebSite, error) {
	website := latencyUseCase.repo.GetMaxLatencyUnit()
	if website == nil || website.Latency == nil {
		return nil, fmt.Errorf("no data")
	}

	return website, nil
}

func (latencyUseCase *LatencyUseCase) GetLatencyLogic(siteName string) (*entity.WebSite, error) {
	website, ok := latencyUseCase.repo.GetLatencyUnit(siteName)

	if (website == nil || website.Latency == nil) && !ok {
		return website, fmt.Errorf("no data")
	}

	return website, nil
}
