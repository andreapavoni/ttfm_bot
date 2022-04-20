package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func EscortCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleAdmin, ttfm.UserRoleBotModerator},
		Help:               "Escort dj off the stage. Usage: `escort <username>`",
		Handler:            escortCommandHandler,
	}
}

func escortCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.Users.UserFromId(cmd.UserId)

	if len(cmd.Args) < 1 {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("You must specify the username of the user you want to escort")}
	}

	escortedUser, err := b.Users.UserFromName(strings.Join(cmd.Args, " "))
	if err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: fmt.Errorf("I can't find the user you want to escort: @%s", escortedUser.Name)}
	}

	if err := b.Actions.EscortDj(escortedUser.Id); err != nil {
		return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypePm, Err: fmt.Errorf("I was unable to escort @%s", escortedUser.Name)}

	}
	return &ttfm.CommandOutput{User: user, ReplyType: ttfm.MessageTypeNone}
}
