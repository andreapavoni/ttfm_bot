package commands

import (
	"errors"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func PlaylistRemoveSongCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleAdmin},
		Help:               "Remove current playing song from current playlist",
		Handler:            playlistRemoveSongCommandHandler,
	}
}

func playlistRemoveSongCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)
	if err := b.RemoveSongFromPlaylist(b.Room.Song.Id); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: errors.New("I was unable to delete the playlist: " + err.Error())}
	}

	return &ttfm.CommandOutput{Msg: "I've removed the song from the current playlist", User: user, ReplyType: ttfm.MessageTypePm}
}
