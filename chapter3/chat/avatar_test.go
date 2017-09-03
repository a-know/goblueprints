package main

import "testing"

func TestAuthAvatarNoURL(t *testing.T) {
	var authAvatar AuthAvatar
	client := new(client)
	_, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("If url not exists, AuthAvatar.GetAvatarURL should return ErrNoAvatarURL.")
	}
}

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	client := new(client)
	url, err := authAvatar.GetAvatarURL(client)
	testUrl := "http://wrl-to-avatar/"
	client.userData = map[string]interface{}{"avatar_url": testUrl}
	url, err = authAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("If url exists, AuthAvatar.GetAvatarURL should not return error.")
	} else {
		if url != testUrl {
			t.Error("AuthAvatar.GetAvatarURL should return correct url.")
		}
	}
}
