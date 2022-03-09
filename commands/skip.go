package commands

import (
	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func SkipCommandHandler(b *ttfm.Bot, userId string, args []string) (string, *ttfm.User, error) {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return "", user, err
	}

	if err := requireBotModerator(b, user); err != nil {
		return "", user, err
	}

	b.SkipSong()

	return "", user, nil
}
