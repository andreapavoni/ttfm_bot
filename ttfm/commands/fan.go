package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func FanCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleAdmin},
		Help:               "Fan user",
		Handler:            fanCommandHandler,
	}
}

func fanCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)

	if len(cmd.Args) < 1 {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("You must specify the username of the user you want to become a fan")}
	}

	fannedUser, err := b.Users.UserFromName(strings.Join(cmd.Args, " "))

	if err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: fmt.Errorf("I can't find the user you want to fan: @%s", fannedUser.Name)}
	}

	if err := b.Users.Fan(fannedUser.Id); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: fmt.Errorf("I was unable to fan @%s", fannedUser.Name)}
	}

	msg := fmt.Sprintf("/me became a fan of @%s", fannedUser.Name)
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}
