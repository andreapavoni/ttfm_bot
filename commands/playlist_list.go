package commands

import (
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func PlaylistListCommandHandler(b *ttfm.Bot, userId string, args []string) (string, *ttfm.User, error) {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return "", user, err
	}

	pls := []string{}

	for _, pl := range b.Playlists.List() {
		if pl == b.Config.CurrentPlaylist {
			pl = "*" + pl
		}
		pls = append(pls, pl)
	}

	msg := "Available playists (the current one is prefixed with a *): " + strings.Join(pls, ", ")

	return msg, user, nil
}
