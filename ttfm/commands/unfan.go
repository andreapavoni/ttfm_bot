package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func UnfanCommand() *ttfm.Command {
	return &ttfm.Command{AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleAdmin}, Help: "Unfan user. Usage: `unfan <username>`", Handler: unfanCommandHandler}
}

func unfanCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)
	if len(cmd.Args) < 1 {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("You must specify the username of the user you want to unfan")}
	}

	unfannedUser, err := b.Users.UserFromName(strings.Join(cmd.Args, " "))
	if err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: errors.New("I can't find the user you want to unfan")}
	}

	if err := b.Users.Unfan(unfannedUser.Id); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: fmt.Errorf("I was unable to unfan @%s", unfannedUser.Name)}
	}

	msg := fmt.Sprintf("/me unfanned @%s", unfannedUser.Name)
	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}
