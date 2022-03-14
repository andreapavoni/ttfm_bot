package commands

import (
	"errors"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func PlaylistRemoveSongCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: err}
	}

	if err := b.RemoveSongFromPlaylist(b.Room.Song.Id); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: errors.New("I was unable to delete the playlist: " + err.Error())}
	}

	return &ttfm.CommandOutput{Msg: "I've removed the song from the current playlist", User: user, ReplyType: ttfm.MessageTypePm}
}
