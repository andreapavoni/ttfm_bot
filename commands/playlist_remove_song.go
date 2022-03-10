package commands

import (
	"errors"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func PlaylistRemoveSongCommandHandler(b *ttfm.Bot, userId string, args []string) (string, *ttfm.User, error) {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return "", user, err
	}

	if err := b.RemoveSongFromPlaylist(b.Room.Song.Id); err != nil {
		return "", user, errors.New("I was unable to delete the playlist: " + err.Error())
	}

	return "I've removed the song from the current playlist", user, nil
}
