package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func BootCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm}

	}

	if err := requireBotModerator(b, user); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm}
	}

	if len(cmd.Args) < 1 {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("You must specify the username of the user you kick")}
	}

	kickedUser, err := b.UserFromName(strings.Join(cmd.Args, " "))
	reason := fmt.Sprintf("Ask @%s for details", user.Name)

	if err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: errors.New("I can't find the user you want to kick")}
	}

	if err := b.BootUser(kickedUser.Id, reason); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: errors.New("I wasn't able to boot @" + kickedUser.Name)}
	}

	return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypeNone}
}
