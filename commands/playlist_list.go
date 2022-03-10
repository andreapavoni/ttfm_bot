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

	msg := "Available playists: " + strings.Join(b.Playlists.List(), ", ")

	return msg, user, nil
}
