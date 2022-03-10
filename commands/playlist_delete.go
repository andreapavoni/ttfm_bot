package commands

import (
	"errors"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func PlaylistDeleteCommandHandler(b *ttfm.Bot, userId string, args []string) (string, *ttfm.User, error) {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return "", user, err
	}

	if len(args) < 1 {
		return "", user, errors.New("You must specify a name of the playlist you want to delete")
	}

	playlistName := strings.Join(args, " ")

	if err := b.RemovePlaylist(playlistName); err != nil {
		return "", user, errors.New("I was unable to delete the playlist: " + err.Error())
	}

	return "I've deleted the playlist", user, nil
}
