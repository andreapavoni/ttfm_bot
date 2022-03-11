package commands

import (
	"errors"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func PlaylistRemoveSongCommandHandler(b *ttfm.Bot, userId string, args []string) *ttfm.CommandOutput {
	user, _ := b.UserFromId(userId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: err}
	}

	if err := b.RemoveSongFromPlaylist(b.Room.Song.Id); err != nil {
		return &ttfm.CommandOutput{Msg: "", User: user, ReplyWith: "pm", Err: errors.New("I was unable to delete the playlist: " + err.Error())}
	}

	return &ttfm.CommandOutput{Msg: "I've removed the song from the current playlist", User: user, ReplyWith: "pm", Err: nil}
}
