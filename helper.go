package gosms_starter

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/kordar/gosms"
)

var (
	providers = make(map[string]gosms.SMSProvider)
	mu        sync.RWMutex
)

func Get(name string) gosms.SMSProvider {
	mu.RLock()
	defer mu.RUnlock()
	p, ok := providers[name]
	if !ok {
		slog.Error("sms provider not exist", "name", name)
		panic(fmt.Errorf("sms provider %s not exist", name))
	}
	return p
}

func Provide(name string, p gosms.SMSProvider) {
	mu.Lock()
	defer mu.Unlock()
	providers[name] = p
}

func ProvideFromConfig(name string, cfg *gosms.SMSConfig) (gosms.SMSProvider, error) {
	p, err := gosms.NewSMSProvider(cfg)
	if err != nil {
		return nil, err
	}
	Provide(name, p)
	return p, nil
}

func ProvideEFromConfig(name string, cfg *gosms.SMSConfig) gosms.SMSProvider {
	p, err := ProvideFromConfig(name, cfg)
	if err != nil {
		slog.Error("provide sms provider failed", "name", name, "err", err)
		panic(fmt.Errorf("provide %s failed: %w", name, err))
	}
	return p
}
