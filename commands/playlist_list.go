package commands

import (
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func PlaylistListCommandHandler(b *ttfm.Bot, userId string, args []string) *ttfm.CommandOutput {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: err}
	}

	pls := []string{}

	for _, pl := range b.Playlists.List() {
		if pl == b.Config.CurrentPlaylist {
			pl = "*" + pl
		}
		pls = append(pls, pl)
	}

	msg := "Available playists (the current one is prefixed with a *): " + strings.Join(pls, ", ")
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyWith: "pm", Err: nil}
}
