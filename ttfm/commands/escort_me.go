package commands

import (
	"fmt"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func EscortMeCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleBotModerator},
		Help:               "Ask to be escorted after current song has been played",
		Handler:            escortMeCommandHandler,
	}
}

func escortMeCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)
	b.AddDjEscorting(cmd.UserId)

	var msg string

	if cmd.Source == ttfm.MessageTypeRoom {
		msg = fmt.Sprintf("@%s I'll escort you at the end of this song", user.Name)
	} else {
		msg = "I'll escort you at the end of this song"
	}

	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}
