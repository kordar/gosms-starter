package gosms_starter

import (
	logger "github.com/kordar/gologger"
	"github.com/kordar/gosms"
	"github.com/spf13/cast"
)

type SMSModule struct {
	name string
	load func(moduleName string, itemId string, p gosms.SMSProvider, item map[string]interface{})
}

func NewSMSModule(name string, load func(moduleName string, itemId string, p gosms.SMSProvider, item map[string]interface{})) *SMSModule {
	return &SMSModule{name: name, load: load}
}

func (m SMSModule) Name() string {
	return m.name
}

func (m SMSModule) _load(id string, cfg map[string]interface{}) {
	if id == "" {
		logger.Fatalf("[%s] the attribute id cannot be empty.", m.Name())
		return
	}

	provider := cast.ToString(cfg["provider"])
	accessKey := cast.ToString(cfg["access_key"])
	secretKey := cast.ToString(cfg["secret_key"])
	signName := cast.ToString(cfg["sign"])
	templateID := cast.ToString(cfg["template"])

	if provider == "" {
		logger.Fatalf("[%s] id=%s provider cannot be empty", m.Name(), id)
		return
	}

	smsCfg := gosms.NewSMSConfig(provider, accessKey, secretKey)
	if signName != "" {
		smsCfg.WithSign(signName)
	}
	if templateID != "" {
		smsCfg.WithTemplate(templateID)
	}

	extra := cast.ToStringMapString(cfg["extra"])
	for k, v := range extra {
		smsCfg.WithExtraParam(k, v)
	}

	p, err := ProvideFromConfig(id, smsCfg)
	if err != nil {
		logger.Fatalf("[%s] id=%s err=%v", m.Name(), id, err)
		return
	}

	m.load(m.Name(), id, p, cfg)
	logger.Infof("[%s] loading module '%s' successfully", m.Name(), id)
}

func (m SMSModule) Load(value interface{}) {
	if value == nil {
		return
	}

	items := cast.ToStringMap(value)
	if items["id"] != nil {
		id := cast.ToString(items["id"])
		m._load(id, items)
		return
	}

	for key, item := range items {
		m._load(key, cast.ToStringMap(item))
	}
}

func (m SMSModule) Close() {
}
