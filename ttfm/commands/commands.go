package commands

import (
	"errors"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func requireAdmin(b *ttfm.Bot, u *ttfm.User) error {
	if !b.UserIsAdmin(u) {
		return errors.New("I won't obey you because you aren't one of my admins or we aren't in the same room")
	}

	return nil
}

func requireBotModerator(b *ttfm.Bot, u *ttfm.User) error {
	if !b.UserIsModerator(u.Id) {
		return errors.New("Sorry, I can't proceed because I'm not a moderator in this room")
	}

	return nil
}
