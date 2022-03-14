package commands

import (
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func PlaylistListCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: err}
	}

	pls := []string{}

	for _, pl := range b.Playlists.List() {
		if pl == b.Config.CurrentPlaylist {
			pl = "[" + pl + "]"
		}
		pls = append(pls, pl)
	}

	msg := "Available playists (the current one is prefixed with a *): " + strings.Join(pls, ", ")
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}
