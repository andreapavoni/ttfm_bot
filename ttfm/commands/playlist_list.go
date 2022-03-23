package commands

import (
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func PlaylistListCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleAdmin},
		Help:               "List playlists",
		Handler:            playlistListCommandHandler,
	}
}

func playlistListCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)
	pls := []string{}

	for _, pl := range b.Playlists.List() {
		if pl == b.Config.CurrentPlaylist {
			pl = "[" + pl + "]"
		}
		pls = append(pls, pl)
	}

	msg := "Available playists (the current one is highlighted): " + strings.Join(pls, ", ")
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}
