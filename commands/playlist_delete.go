package commands

import (
	"errors"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func PlaylistDeleteCommandHandler(b *ttfm.Bot, userId string, args []string) *ttfm.CommandOutput {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: err}
	}

	if len(args) < 1 {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: errors.New("You must specify a name of the playlist you want to delete")}
	}

	playlistName := strings.Join(args, " ")

	if err := b.RemovePlaylist(playlistName); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: errors.New("I was unable to delete the playlist: " + err.Error())}
	}

	return &ttfm.CommandOutput{Msg: "I've deleted the playlist", User: user, ReplyWith: "pm", Err: nil}
}
