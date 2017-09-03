package main

import (
	"errors"
)

// Error occurs in can not return avatar url
var ErrNoAvatarURL = errors.New("chat: can not get avatar url.")

type Avatar interface {
	// return avatar URL.
	// return Error when can not get avatar url.
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}
