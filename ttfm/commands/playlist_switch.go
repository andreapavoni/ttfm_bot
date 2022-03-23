package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func PlaylistSwitchCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleAdmin},
		Help:               "Switch to another playlist",
		Handler:            playlistSwitchCommandHandler,
	}
}

func playlistSwitchCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)
	if len(cmd.Args) < 1 {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("You must specify a name of the playlist you want to switch to")}
	}

	playlistName := strings.Join(cmd.Args, " ")

	if err := b.Playlists.Switch(playlistName); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: fmt.Errorf("I was unable to switch playlist: %s", err.Error())}
	}

	msg := fmt.Sprintf("/me switched to playlist `%s`", playlistName)
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}
