package discord

import (
	"testing"

	"golang.org/x/oauth2"
)

func TestTokenNil(t *testing.T) {
	r := convertToken(nil)
	if r != nil {
		t.Error("Should be nil")
	}
}

func TestNotNil(t *testing.T) {
	tk := oauth2.Token{}

	r := convertToken(&tk)
	if r == nil {
		t.Error("Should not be nil")
	}
}

func TestExtaMap(t *testing.T) {
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

	r := convertToken(tk)
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
