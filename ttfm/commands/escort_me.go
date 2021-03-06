package commands

import (
	"fmt"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func EscortMeCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleBotModerator},
		Help:               "Ask to be escorted after your song has been played",
		Handler:            escortMeCommandHandler,
	}
}

func escortMeCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)
	if err := b.Actions.AddDjEscorting(cmd.UserId); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: err}
	}

	var msg string

	if cmd.Source == ttfm.MessageTypeRoom {
		msg = fmt.Sprintf("@%s I'll escort you at the end of your played song", user.Name)
	} else {
		msg = "I'll escort you at the end of your played song"
	}

	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}
