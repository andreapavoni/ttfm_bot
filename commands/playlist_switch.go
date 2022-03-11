package commands

import (
	"errors"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func PlaylistSwitchCommandHandler(b *ttfm.Bot, userId string, args []string) *ttfm.CommandOutput {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: err}
	}

	if len(args) < 1 {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: errors.New("You must specify a name of the playlist you want to switch to")}
	}

	playlistName := strings.Join(args, " ")

	if err := b.SwitchPlaylist(playlistName); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: errors.New("I was unable to switch playlist: " + err.Error())}
	}

	return &ttfm.CommandOutput{Msg: "/me switched playlist", User: user, ReplyWith: "action", Err: nil}
}
