package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func FanCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)

	if err := requireAdmin(b, user); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: err}
	}

	if len(cmd.Args) < 1 {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("You must specify the username of the user you want to become a fan")}
	}

	fannedUser, err := b.UserFromName(strings.Join(cmd.Args, " "))

	if err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: errors.New("I can't find the user you want to fan: @" + fannedUser.Name)}
	}

	if err := b.Fan(fannedUser.Id); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: errors.New("I was unable to fan @" + fannedUser.Name)}
	}

	msg := fmt.Sprintf("/me became a fan of @%s", fannedUser.Name)
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}
