package commands

import (
	"errors"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func EscortCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: err}
	}

	if err := requireBotModerator(b, user); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: err}
	}

	if len(cmd.Args) < 1 {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: errors.New("You must specify the username of the user you want to escort")}
	}

	escortedUser, err := b.UserFromName(strings.Join(cmd.Args, " "))
	if err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: errors.New("I can't find the user you want to escort: @" + escortedUser.Name)}
	}

	if err := b.EscortDj(escortedUser.Id); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: errors.New("I was unable to escort @" + escortedUser.Name)}

	}

	return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypeNone}
}
