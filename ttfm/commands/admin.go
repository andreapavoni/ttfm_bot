package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func AdminCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleAdmin, ttfm.UserRoleBotModerator},
		Help:               "Add or remove bot admin. Usage: `admin [<add|remove> <username>]`. Without arguments it shows the current admins",
		Handler:            adminCommandHandler,
	}
}

func adminCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)

	if len(cmd.Args) == 0 {
		msg := "Current admins: " + strings.Join(b.Admins.Values(), ", ")
		return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
	}

	if len(cmd.Args) < 2 {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("You must specify the action (add or remove) and username of the user")}
	}

	userObj, err := b.Users.UserFromName(cmd.Args[1])
	if err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: fmt.Errorf("I can't find the user called `%s`", cmd.Args[1])}
	}

	var msg string
	switch cmd.Args[0] {
	case "add":
		if err := b.Admins.Put(userObj.Id, userObj.Name); err != nil {
			return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: fmt.Errorf("I was unable to add `@%s` to admins: %s", userObj.Name, err.Error())}
		}
		msg = fmt.Sprintf("Added @%s to my admins", userObj.Name)
	case "remove":
		if err := b.Admins.Delete(userObj.Id); err != nil {
			return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: fmt.Errorf("I was unable to remove `@%s` from admins: %s", userObj.Name, err.Error())}
		}
		msg = fmt.Sprintf("Removed @%s from my admins", userObj.Name)
	default:
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: fmt.Errorf("You should use `add` or `remove` as action")}
	}

	return &ttfm.CommandOutput{Msg: msg, User: user, ReplyType: cmd.Source}
}
