package httpbuilder

import (
	"context"
	"net/http"
	"os"
	"testing"
)

// DiscordUser for testing
type DiscordUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
}

func TestGetDiscordUser(t *testing.T) {
	const url = `https://discordapp.com/api/users/@me`
	var botToken = os.Getenv("DISCORD_BOT_TOKEN")

	ctx := context.Background()

	var result DiscordUser
	builder := New().
		SetURL(url).
		SetOut(&result).
		SetAuthToken("Bot", botToken).
		SetLogger(t.Logf)

	status, err := builder.Do(ctx)
	if err != nil {
		t.Error(err)
	}

	if status != http.StatusOK {
		t.Errorf("Wrong status: %d", status)
	}

	// we should get a result!
	if result.ID != "492157091893477377" {
		t.Errorf("Wrong user id: %s", result.ID)
	}

	if result.Username != "ctx" {
		t.Errorf("Wrong username: %s", result.Username)
	}
}
