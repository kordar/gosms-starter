package gosms_starter

import (
	"fmt"
	"log/slog"
	"strings"

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
		slog.Error("sms module id cannot be empty", "module", m.Name())
		panic(fmt.Errorf("[%s] the attribute id cannot be empty", m.Name()))
	}

	provider := cast.ToString(cfg["provider"])
	accessKey := cast.ToString(cfg["access_key"])
	secretKey := cast.ToString(cfg["secret_key"])
	signName := cast.ToString(cfg["sign"])
	templateID := cast.ToString(cfg["template"])

	if provider == "" {
		slog.Error("sms provider cannot be empty", "module", m.Name(), "id", id)
		panic(fmt.Errorf("[%s] id=%s provider cannot be empty", m.Name(), id))
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

	extrastr := cast.ToString(cfg["extrastr"])
	if extrastr != "" {
		for _, item := range strings.Split(extrastr, ",") {
			if item == "" {
				continue
			}
			kv := strings.Split(item, "::")
			if len(kv) != 2 {
				continue
			}
			smsCfg.WithExtraParam(kv[0], kv[1])
		}
	}

	p, err := ProvideFromConfig(id, smsCfg)
	if err != nil {
		slog.Error("provide sms provider failed", "module", m.Name(), "id", id, "err", err)
		panic(fmt.Errorf("[%s] id=%s err=%w", m.Name(), id, err))
	}

	if m.load != nil {
		m.load(m.Name(), id, p, cfg)
	}
	slog.Info("sms module loaded", "module", m.Name(), "id", id)
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
