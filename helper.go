package gosms_starter

import (
	"sync"

	logger "github.com/kordar/gologger"
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
		logger.Fatalf("sms provider %s not exist.", name)
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
		logger.Fatalf("[provide %s] %v", name, err)
	}
	return p
}
