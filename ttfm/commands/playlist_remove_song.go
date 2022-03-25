package commands

import (
	"fmt"

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
	user, _ := b.Users.UserFromId(cmd.UserId)
	if err := b.CurrentPlaylist.RemoveSong(b.Room.Song.Id); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: fmt.Errorf("I was unable to delete the playlist: %s", err.Error())}
	}

	return &ttfm.CommandOutput{Msg: "/me removed the song from the current playlist", User: user, ReplyType: cmd.Source}
}
