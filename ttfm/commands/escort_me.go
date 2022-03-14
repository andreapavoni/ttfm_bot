package commands

import (
	"fmt"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func EscortMeCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)

	if err := requireBotModerator(b, user); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: err}
	}

	b.AddDjEscorting(cmd.UserId)

	var msg string

	if cmd.Source == ttfm.MessageTypeRoom {
		msg = fmt.Sprintf("@%s I'll escort you at the end of this song", user.Name)
	} else {
		msg = "I'll escort you at the end of this song"
	}

	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}
