package smashgg

import "testing"

func TestExtract(t *testing.T) {
	userURL := "https://smash.gg/admin/user/46ee4f62/profile-settings"

	slug, err := extractSlug(userURL)
	if err != nil {
		t.Error(err)
		return
	}

	if slug != "user/46ee4f62" {
		t.Errorf("Invalid user slug %s", slug)
		return
	}
}
