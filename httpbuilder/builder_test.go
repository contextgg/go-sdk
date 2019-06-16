package httpbuilder

import (
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
	const url = "https://discordapp.com/api/users/234466834202492928"
	var botToken = os.Getenv("DISCORD_BOT_TOKEN")

	var result DiscordUser
	builder := New().
		SetURL(url).
		SetOut(&result).
		SetAuthToken("Bot", botToken).
		SetLogger(t.Logf)

	status, err := builder.Do()
	if err != nil {
		t.Error(err)
	}

	if status != http.StatusOK {
		t.Errorf("Wrong status: %d", status)
	}

	// we should get a result!
	if result.ID != "234466834202492928" {
		t.Errorf("Wrong user id: %s", result.ID)
	}

	if result.Username != "Doofus Viper" {
		t.Errorf("Wrong username: %s", result.Username)
	}
}
