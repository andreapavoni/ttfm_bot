package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func PlaylistSwitchCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: err}
	}

	if len(cmd.Args) < 1 {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("You must specify a name of the playlist you want to switch to")}
	}

	playlistName := strings.Join(cmd.Args, " ")

	if err := b.SwitchPlaylist(playlistName); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: errors.New("I was unable to switch playlist: " + err.Error())}
	}

	msg := fmt.Sprintf("/me switched to playlist `%s`", playlistName)
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}
