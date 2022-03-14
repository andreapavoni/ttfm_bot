package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func UnfanCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: err}
	}

	if len(cmd.Args) < 1 {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("You must specify the username of the user you want to unfan")}
	}

	unfannedUser, err := b.UserFromName(strings.Join(cmd.Args, " "))
	if err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: errors.New("I can't find the user you want to unfan")}
	}

	if err := b.Unfan(unfannedUser.Id); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: errors.New("I was unable to unfan @" + unfannedUser.Name)}
	}

	msg := fmt.Sprintf("/me unfanned @%s", unfannedUser.Name)
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}
