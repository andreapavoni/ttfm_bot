package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func PlaylistDeleteCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleAdmin},
		Help:               "Delete playlist",
		Handler:            playlistDeleteCommandHandler,
	}
}

func playlistDeleteCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)

	if len(cmd.Args) < 1 {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("You must specify a name of the playlist you want to delete")}
	}

	playlistName := strings.Join(cmd.Args, " ")

	if err := b.Playlists.Remove(playlistName); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: fmt.Errorf("I was unable to delete the playlist: %s", err.Error())}
	}

	msg := fmt.Sprintf("/me deleted playlist `%s`", playlistName)
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}
