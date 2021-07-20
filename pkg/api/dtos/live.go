package dtos

import (
	"encoding/json"

	"github.com/grafana/grafana/pkg/models"
)

type LivePublishCmd struct {
	Channel string          `json:"channel"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type LivePublishResponse struct {
}

type LiveChannelRule struct {
	Id               int64                             `json:"id"`
	Version          int                               `json:"version"`
	Pattern          string                            `json:"pattern"`
	Config           models.LiveChannelRulePlainConfig `json:"config"`
	SecureJsonFields map[string]bool                   `json:"secureJsonFields"`
}

type LiveChannelRuleListItem struct {
	Id      int64  `json:"id"`
	Version int    `json:"version"`
	Pattern string `json:"pattern"`
}
