package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func PlaylistCreateCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleAdmin},
		Help:               "Create new playlist",
		Handler:            playlistCreateCommandHandler,
	}
}

func playlistCreateCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)
	if len(cmd.Args) < 1 {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("You must specify a name of the new playlist")}
	}

	playlistName := strings.Join(cmd.Args, " ")

	if err := b.AddPlaylist(playlistName); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: errors.New("I was unable to add the new playlist: " + err.Error())}
	}

	msg := fmt.Sprintf("/me created playlist `%s`", playlistName)
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}
