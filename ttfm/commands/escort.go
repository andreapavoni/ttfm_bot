package commands

import (
	"errors"
	"strings"

	"github.com/andreapavoni/ttfm_bot/ttfm"
)

func EscortCommand() *ttfm.Command {
	return &ttfm.Command{
		AuthorizationRoles: []ttfm.UserRole{ttfm.UserRoleAdmin, ttfm.UserRoleBotModerator},
		Help:               "Escort dj off the stage",
		Handler:            escortCommandHandler,
	}
}

func escortCommandHandler(b *ttfm.Bot, cmd *ttfm.CommandInput) *ttfm.CommandOutput {
	user, _ := b.UserFromId(cmd.UserId)

	if len(cmd.Args) < 1 {
		return &ttfm.CommandOutput{User: user, ReplyType: cmd.Source, Err: errors.New("You must specify the username of the user you want to escort")}
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
