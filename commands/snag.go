package commands

import (
	"errors"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func SnagCommandHandler(b *ttfm.Bot, userId string, args []string) (string, *ttfm.User, error) {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return "", user, err
	}

	if err := b.Snag(b.Room.Song.Id); err == nil {
		return "I snagged this song!", user, nil
	}

	return "", user, errors.New("I've failed to snag this song")
}
