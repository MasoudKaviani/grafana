package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/grafana/grafana/pkg/api/dtos"
	"github.com/grafana/grafana/pkg/api/response"
	"github.com/grafana/grafana/pkg/models"
	"github.com/grafana/grafana/pkg/util"
)

//go:generate mockgen -destination=live_channel_rule_mock.go -package=api github.com/grafana/grafana/pkg/api ChannelRuleStorage

type ChannelRuleStorage interface {
	ListChannelRules(ctx context.Context, cmd models.ListLiveChannelRuleCommand) ([]*models.LiveChannelRule, error)
	GetChannelRule(ctx context.Context, cmd models.GetLiveChannelRuleCommand) (*models.LiveChannelRule, error)
	CreateChannelRule(ctx context.Context, cmd models.CreateLiveChannelRuleCommand) (*models.LiveChannelRule, error)
	UpdateChannelRule(ctx context.Context, cmd models.UpdateLiveChannelRuleCommand) (*models.LiveChannelRule, error)
	DeleteChannelRule(ctx context.Context, cmd models.DeleteLiveChannelRuleCommand) (int64, error)
}

type channelRuleAPI struct {
	storage ChannelRuleStorage
}

func (a *channelRuleAPI) ListChannelRules(c *models.ReqContext) response.Response {
	query := models.ListLiveChannelRuleCommand{OrgId: c.OrgId}

	result, err := a.storage.ListChannelRules(c.Req.Context(), query)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "Failed to query channel rules", err)
	}

	items := make([]dtos.LiveChannelRuleListItem, 0, len(result))

	for _, ch := range result {
		item := dtos.LiveChannelRuleListItem{
			Id:      ch.Id,
			Version: ch.Version,
			Pattern: ch.Pattern,
		}
		items = append(items, item)
	}
	return response.JSON(http.StatusOK, &items)
}

func liveChannelToDTO(ch *models.LiveChannelRule) dtos.LiveChannelRule {
	item := dtos.LiveChannelRule{
		Id:               ch.Id,
		Version:          ch.Version,
		Pattern:          ch.Pattern,
		Config:           ch.Config,
		SecureJsonFields: map[string]bool{},
	}
	for k, v := range ch.Secure {
		if len(v) > 0 {
			item.SecureJsonFields[k] = true
		}
	}
	return item
}

func (a *channelRuleAPI) GetChannelRuleById(c *models.ReqContext) response.Response {
	query := models.GetLiveChannelRuleCommand{
		Id:    c.ParamsInt64(":id"),
		OrgId: c.OrgId,
	}

	result, err := a.storage.GetChannelRule(c.Req.Context(), query)
	if err != nil {
		if errors.Is(err, models.ErrLiveChannelRuleNotFound) {
			return response.Error(http.StatusNotFound, "Channel rule not found", nil)
		}
		return response.Error(http.StatusInternalServerError, "Failed to query channel rule", err)
	}
	item := liveChannelToDTO(result)
	return response.JSON(http.StatusOK, &item)
}

func (a *channelRuleAPI) CreateChannelRule(c *models.ReqContext, cmd models.CreateLiveChannelRuleCommand) response.Response {
	cmd.OrgId = c.OrgId

	result, err := a.storage.CreateChannelRule(c.Req.Context(), cmd)
	if err != nil {
		if errors.Is(err, models.ErrLiveChannelRuleExists) {
			return response.Error(http.StatusConflict, err.Error(), err)
		}
		return response.Error(http.StatusInternalServerError, "Failed to create channel rule", err)
	}

	return response.JSON(http.StatusOK, util.DynMap{
		"message":     "channel rule added",
		"id":          result.Id,
		"channelRule": liveChannelToDTO(result),
	})
}

func (a *channelRuleAPI) fillChannelRuleWithSecureJSONData(ctx context.Context, cmd *models.UpdateLiveChannelRuleCommand) error {
	if len(cmd.Secure) == 0 {
		return nil
	}

	rule, err := a.storage.GetChannelRule(ctx, models.GetLiveChannelRuleCommand{
		OrgId: cmd.OrgId,
		Id:    cmd.Id,
	})
	if err != nil {
		return err
	}

	secureJSONData := rule.Secure.Decrypt()
	for k, v := range secureJSONData {
		if _, ok := cmd.Secure[k]; !ok {
			cmd.Secure[k] = v
		}
	}

	return nil
}

func (a *channelRuleAPI) UpdateChannelRule(c *models.ReqContext, cmd models.UpdateLiveChannelRuleCommand) response.Response {
	cmd.Id = c.ParamsInt64(":id")
	cmd.OrgId = c.OrgId

	err := a.fillChannelRuleWithSecureJSONData(c.Req.Context(), &cmd)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "Failed to update channel rule", err)
	}

	_, err = a.storage.UpdateChannelRule(c.Req.Context(), cmd)
	if err != nil {
		if errors.Is(err, models.ErrLiveChannelRuleNotFound) {
			return response.Error(http.StatusNotFound, "Channel rule not found", nil)
		}
		if errors.Is(err, models.ErrLiveChannelRuleUpdatingOldVersion) {
			return response.Error(http.StatusConflict, "Channel rule has already been updated by someone else. Please reload and try again", err)
		}
		return response.Error(http.StatusInternalServerError, "Failed to update channel rule", err)
	}

	getCmd := models.GetLiveChannelRuleCommand{
		Id:    cmd.Id,
		OrgId: c.OrgId,
	}

	result, err := a.storage.GetChannelRule(c.Req.Context(), getCmd)
	if err != nil {
		if errors.Is(err, models.ErrLiveChannelRuleNotFound) {
			return response.Error(http.StatusNotFound, "Channel rule not found", nil)
		}
		return response.Error(http.StatusInternalServerError, "Failed to query channel rule", err)
	}

	return response.JSON(http.StatusOK, util.DynMap{
		"message":     "channel rule updated",
		"id":          cmd.Id,
		"channelRule": liveChannelToDTO(result),
	})
}

func (a *channelRuleAPI) DeleteChannelRuleById(c *models.ReqContext) response.Response {
	id := c.ParamsInt64(":id")

	if id <= 0 {
		return response.Error(http.StatusBadRequest, "Missing valid channel rule id", nil)
	}

	getCmd := models.GetLiveChannelRuleCommand{
		Id:    id,
		OrgId: c.OrgId,
	}
	_, err := a.storage.GetChannelRule(c.Req.Context(), getCmd)
	if err != nil {
		if errors.Is(err, models.ErrLiveChannelRuleNotFound) {
			return response.Error(http.StatusNotFound, "Channel rule not found", nil)
		}
		return response.Error(http.StatusInternalServerError, "Failed to query channel rule", err)
	}

	cmd := models.DeleteLiveChannelRuleCommand{Id: id, OrgId: c.OrgId}
	_, err = a.storage.DeleteChannelRule(c.Req.Context(), cmd)
	if err != nil {
		return response.Error(http.StatusInternalServerError, "Failed to delete channel rule", err)
	}
	return response.Success("Channel rule deleted")
}
