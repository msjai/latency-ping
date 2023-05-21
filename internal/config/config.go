package config

import (
	"bufio"
	"errors"
	"flag"
	"os"
	"strconv"

	"go.uber.org/zap"

	"github.com/msjai/latency-ping/internal/logger"
)

var (
	ErrNoRunAddress  = errors.New("config - New - RunAddressFromCMA: run address not set in env and cla")
	ErrNoWorkerCount = errors.New("config - New - WorkerCountFromCMA: worker count not set in env and cla")
)

// Config -.
type Config struct {
	RunAddress   string   // Адрес и порт запуска сервиса
	ListWebSites []string // слайс частей адреса сайта, это названия сайта без https://
	WorkerCount  int
	L            *zap.SugaredLogger // Логгер// Адрес и порт запуска сервиса
}

// New returns app config
func New() (*Config, error) {
	cfg := &Config{}

	cfg.L = logger.New()

	cfg.RunAddress = os.Getenv("RUN_ADDRESS")
	if cfg.RunAddress == "" {
		cfg.L.Infow("config info: server address not set in env")
		RunAddressFromCMA(cfg)
		if cfg.RunAddress == "" {
			cfg.L.Infow("config info: server address not set in cla")
			return cfg, ErrNoRunAddress
		}
	}

	workerCountString, _ := strconv.Atoi(os.Getenv("WORKER_COUNT"))
	cfg.WorkerCount = workerCountString
	if cfg.WorkerCount == 0 {
		cfg.L.Infow("config info: worker count not set in env")
		RunWorkerCountFromCMA(cfg)
		if cfg.WorkerCount == 0 {
			cfg.L.Infow("config info: worker count not set in env")
			return cfg, ErrNoWorkerCount
		}
	}

	flag.Parse()

	var err error
	cfg.ListWebSites, err = initListWebSites()
	if err != nil {
		cfg.L.Errorw("Cant get list of websites to ping")
		return cfg, err
	}

	return cfg, nil
}

// RunAddressFromCMA -
func RunAddressFromCMA(cfg *Config) {
	flag.StringVar(&cfg.RunAddress, "a", "", "host(server address) to listen on")
}

func RunWorkerCountFromCMA(cfg *Config) {
	wc := flag.Int("wc", 2, "worker count to get stats")
	cfg.WorkerCount = *wc
}

func initListWebSites() ([]string, error) {
	var result []string

	file, err := os.Open("./list_sites")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
