package commands

import (
	"errors"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func PlaylistCreateCommandHandler(b *ttfm.Bot, userId string, args []string) *ttfm.CommandOutput {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: err}
	}

	if len(args) < 1 {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: errors.New("You must specify a name of the new playlist")}
	}

	playlistName := strings.Join(args, " ")

	if err := b.AddPlaylist(playlistName); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: errors.New("I was unable to add the new playlist: " + err.Error())}
	}

	return &ttfm.CommandOutput{Msg: "I've created the new playlist", User: user, ReplyWith: "pm", Err: nil}
}
