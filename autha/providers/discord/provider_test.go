package discord

import (
	"testing"

	"golang.org/x/oauth2"
)

type NopParams struct {
	data map[string]string
}

func (n *NopParams) Get(k string) string {
	if r, ok := n.data[k]; ok {
		return r
	}
	return ""
}

func TestTokenNil(t *testing.T) {
	r := convertToken(nil, &NopParams{})
	if r != nil {
		t.Error("Should be nil")
	}
}

func TestNotNil(t *testing.T) {
	tk := oauth2.Token{}

	r := convertToken(&tk, &NopParams{})
	if r == nil {
		t.Error("Should not be nil")
	}
}
func TestHasGuild(t *testing.T) {
	tk := oauth2.Token{}
	params := &NopParams{
		data: map[string]string{
			"guild_id":    "yes",
			"permissions": "234",
		},
	}

	r := convertToken(&tk, params)
	if r == nil {
		t.Error("Should not be nil")
		return
	}

	if r.GuildID != "yes" {
		t.Errorf("GuildID wrong; wanted %s got %s", "yes", r.GuildID)
		return
	}
	if r.Permissions != "234" {
		t.Errorf("GuildID wrong; wanted %s got %s", "234", r.Permissions)
		return
	}
}

func TestExtaMapWebhook(t *testing.T) {
	hook := map[string]interface{}{
		"id":         "id",
		"token":      "token",
		"name":       "name",
		"channel_id": "channel_id",
		"guild_id":   "guild_id",
		"avatar":     "avatar",
		"type":       1,
		"url":        "url",
	}
	wrap := map[string]interface{}{
		"webhook": hook,
	}

	tk := &oauth2.Token{}
	tk = tk.WithExtra(wrap)

	r := convertToken(tk, &NopParams{})
	if r == nil {
		t.Error("Should not be nil")
		return
	}

	if r.Webhook == nil {
		t.Error("Webhook should not be nil")
		return
	}

	if r.Webhook.ID != "id" {
		t.Error("Webhook id wrong")
		return
	}

	if r.Webhook.ChannelID != "channel_id" {
		t.Error("Webhook ChannelID wrong")
		return
	}

	if r.Webhook.GuildID != "guild_id" {
		t.Error("Webhook GuildID wrong")
		return
	}
}

func TestExtaMapGuild(t *testing.T) {
	guild := map[string]interface{}{
		"id":       "id",
		"name":     "name",
		"icon":     "icon",
		"owner_id": "owner_id",
	}
	wrap := map[string]interface{}{
		"guild": guild,
	}

	tk := &oauth2.Token{}
	tk = tk.WithExtra(wrap)

	r := convertToken(tk, &NopParams{})
	if r == nil {
		t.Error("Should not be nil")
		return
	}

	if r.Error != nil {
		t.Error(r.Error)
		return
	}

	if r.Guild == nil {
		t.Error("Guild should not be nil")
		return
	}

	if r.Guild.Name != "name" {
		t.Error("Guild Name wrong")
		return
	}

	if r.Guild.ID != "id" {
		t.Error("Guild id wrong")
		return
	}

	if r.Guild.OwnerID != "owner_id" {
		t.Error("Guild OwnerID wrong")
		return
	}

	if r.Guild.Icon != "icon" {
		t.Error("Guild Icon wrong")
		return
	}
}
